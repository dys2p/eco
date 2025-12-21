// Package ssg creates both static and dynamic sites from translated HTML templates and markdown files.
//
// The content directory root may contain:
//
//   - html template files
//   - one folder for each html page, containing markdown files whose filename root is the language prefix, like "en.md"
//   - files and folders which are copied verbatim (see Keep)
//
// The output is like "/en/page.html".
//
// Note that "gotext update" requires a Go module and package for merging translations, accessing message.DefaultCatalog and writing catalog.go.
// While gotext-update-templates has been extended to accept additional directories, a root module and package is still required for static site generation.
//
// For symlink support see [Handler] and [WriteFiles]. Because it partly follows symlinks, you should use this package on trusted input only.
package ssg

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	paths "path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dys2p/eco/lang"
	"gitlab.com/golang-commonmark/markdown"
	"golang.org/x/text/language"
	_ "golang.org/x/text/message"
)

var Keep = []string{
	"ads.txt",
	"app-ads.txt",
	"assets",
	"files",
	"images",
	"sites",
	"static",
}

var md = markdown.New(markdown.HTML(true), markdown.Linkify(false))

// LangOption should be used in templates.
type LangOption struct {
	BCP47    string
	Name     string
	Prefix   string
	Selected bool
}

// SelectLanguage returns a [Language] slice. If if only one language is present, the slice will be empty.
func LangOptions(langs lang.Languages, selected lang.Lang) []LangOption {
	var languages []LangOption
	if len(langs) > 1 {
		for _, l := range langs {
			languages = append(languages, LangOption{
				BCP47:    l.BCP47,
				Name:     strings.ToUpper(l.Prefix),
				Prefix:   l.Prefix,
				Selected: l.Prefix == selected.Prefix,
			})
		}
	}
	return languages
}

type TemplateData struct {
	lang.Lang
	Languages []LangOption // usually empty if only one language is defined
	Onion     bool
	Path      string // without language prefix, for language buttons and hreflang
	Title     string // for <title>
}

func MakeTemplateData(langs lang.Languages, r *http.Request) TemplateData {
	l, path, _ := langs.FromPath(r.URL.Path)
	return TemplateData{
		Lang:      l,
		Languages: LangOptions(langs, l),
		Onion:     strings.HasSuffix(r.Host, ".onion") || strings.Contains(r.Host, ".onion:"),
		Path:      path,
	}
}

// Hreflangs returns <link hreflang> elements for every td.Language, including the selected language.
//
// See also: https://developers.google.com/search/blog/2011/12/new-markup-for-multilingual-content
func (td TemplateData) Hreflangs() template.HTML {
	var b strings.Builder
	for _, l := range td.Languages {
		b.WriteString(fmt.Sprintf(`<link rel="alternate" hreflang="%s" href="/%s/%s">`, l.BCP47, l.Prefix, td.Path))
		b.WriteString("\n")
	}
	return template.HTML(b.String())
}

type Website struct {
	Fsys    fs.FS               // consider wrapping httputil.ModTimeFS around it
	Dynamic map[string]struct { // url path
		Template *template.Template
		Data     TemplateData
	}
	Static []string // url and filesystem paths
}

func MakeWebsite(fsys fs.FS, add *template.Template, langs lang.Languages) (*Website, error) {
	var dynamic = make(map[string]struct {
		Template *template.Template
		Data     TemplateData
	})
	var static []string

	// collect static content and sites
	var sites []string
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("reading root dir: %w", err)
	}
	for _, entry := range entries {
		// follow symlink
		var isDir = entry.IsDir()
		if entry.Type()&fs.ModeSymlink != 0 {
			// get symlink target FileInfo with fs.Stat
			info, err := fs.Stat(fsys, entry.Name())
			if err == nil {
				if info.Mode()&fs.ModeDir != 0 {
					isDir = true
				}
			}
		}

		switch {
		case strings.HasPrefix(entry.Name(), "."):
			continue
		case slices.Contains(Keep, entry.Name()):
			static = append(static, entry.Name())
		case isDir:
			sites = append(sites, entry.Name())
		}
	}

	// prepare site template
	tmpl, err := template.ParseFS(fsys, "*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	if add != nil {
		for _, t := range add.Templates() {
			if t.Tree != nil { // that's possible
				tmpl, err = tmpl.AddParseTree(t.Name(), t.Tree)
				if err != nil {
					return nil, fmt.Errorf("adding additional template %s: %w", t.Name(), err)
				}
			}
		}
	}

	// translate sites
	for _, site := range sites {
		// read markdown files
		var bcp47 []string
		var title []string   // same indices
		var content []string // same indices
		entries, err := fs.ReadDir(fsys, site)
		if err != nil {
			return nil, fmt.Errorf("reading dir %s: %w", site, err)
		}
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			if entry.IsDir() {
				continue
			}
			ext := filepath.Ext(entry.Name())
			root := strings.TrimSuffix(entry.Name(), ext)
			if ext == ".html" || ext == ".md" {
				var filetitle string
				filecontent, err := fs.ReadFile(fsys, filepath.Join(site, entry.Name()))
				if err != nil {
					return nil, fmt.Errorf("reading file: %w", err)
				}
				if ext == ".md" {
					filetitle = getTitleFromMarkdown(string(filecontent))
					filecontent = []byte(md.RenderToString(filecontent))
				}
				bcp47 = append(bcp47, root)
				title = append(title, string(filetitle))
				content = append(content, string(filecontent))
			}
		}
		if len(content) == 0 {
			continue
		}

		// make matcher for available translations
		var haveTags []language.Tag
		for _, have := range bcp47 {
			haveTag, err := language.Parse(have)
			if err != nil {
				return nil, fmt.Errorf("parsing language %s: %w", have, err)
			}
			haveTags = append(haveTags, haveTag)
		}
		matcher := language.NewMatcher(haveTags)

		// assemble site template
		for _, lang := range langs {
			_, index, _ := matcher.Match(lang.Tag)
			tt, err := tmpl.Clone()
			if err != nil {
				return nil, fmt.Errorf("cloning template: %w", err)
			}
			tt, err = tt.Parse(`{{define "content"}}` + content[index] + `{{end}}`) // or parse content into t and then call AddParseTree(content, t.Tree)
			if err != nil {
				return nil, fmt.Errorf("adding content of %s: %w", site, err)
			}
			outpath := filepath.Join(lang.Prefix, site+".html")
			data := MakeTemplateData(langs, httptest.NewRequest(http.MethodGet, "/"+outpath, nil))
			data.Title = title[index]
			// data.Onion is not known yet

			dynamic[outpath] = struct {
				Template *template.Template
				Data     TemplateData
			}{
				Template: tt,
				Data:     data,
			}
		}
	}

	return &Website{
		Fsys:    fsys,
		Dynamic: dynamic,
		Static:  static,
	}, nil
}

// Handler returns a HTTP handler which serves content from fsys.
// It optionally accepts an additional HTML template and a function which makes custom template data.
// For compatibility with WriteFiles, the custom template data struct should embed TemplateData.
//
// Note that embed.FS does not support symlinks. If you use symlinks to share content,
// consider building a go:generate workflow which calls "cp --dereference".
func (ws Website) Handler(makeTemplateData func(*http.Request, TemplateData) any, next http.Handler) http.Handler {
	handler := http.NewServeMux()

	for path, dynamic := range ws.Dynamic {
		path = paths.Join("/", path)
		handler.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			dynamic.Data.Onion = strings.HasSuffix(r.Host, ".onion") || strings.Contains(r.Host, ".onion:")
			var data any
			if makeTemplateData != nil {
				data = makeTemplateData(r, dynamic.Data)
			} else {
				data = dynamic.Data
			}

			if err := dynamic.Template.ExecuteTemplate(w, "html", data); err != nil {
				log.Printf("error executing ssg template %s: %v", path, err)
			}
		})
	}

	for _, path := range ws.Static {
		pattern := paths.Join("/", path) + "/" // trailing slash means prefix match
		handler.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, ws.Fsys, r.URL.Path) // works for dirs and files
		})
	}

	if next != nil {
		handler.Handle("/", next)
	}

	return handler
}

// WriteFiles creates static HTML files. Templates are executed with TemplateData. Symlinks are dereferenced.
func (ws Website) WriteFiles(outDir string, onion bool) error {
	if realOutDir, err := filepath.EvalSymlinks(outDir); err == nil {
		outDir = realOutDir
	}
	if !strings.HasPrefix(outDir, "/tmp/") {
		return errors.New("refusing to write outside of /tmp")
	}
	_ = os.RemoveAll(outDir)

	for path, dynamic := range ws.Dynamic {
		dst := filepath.Join(outDir, path)
		if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
			return fmt.Errorf("error making folder %s: %v", filepath.Dir(dst), err)
		}
		outfile, err := os.Create(dst)
		if err != nil {
			return fmt.Errorf("error opening outfile %s: %v", dst, err)
		}
		defer outfile.Close()

		dynamic.Data.Onion = onion
		err = dynamic.Template.ExecuteTemplate(outfile, "html", dynamic.Data)
		if err != nil {
			return fmt.Errorf("error executing template for %s: %v%s", dst, err, dynamic.Template.DefinedTemplates())
		}
	}

	for _, path := range ws.Static {
		if err := CopyFS(outDir, ws.Fsys, path); err != nil {
			return fmt.Errorf("error copying %s to %s: %v", path, outDir, err)
		}
	}
	return nil
}

func getTitleFromMarkdown(filecontent string) string {
	filecontent = strings.TrimSpace(filecontent)
	firstLine, _, _ := strings.Cut(filecontent, "\n")
	if title, ok := strings.CutPrefix(firstLine, "# "); ok {
		return title
	}
	return ""
}

// ListenAndServe provides an easy way to preview a static site with absolute src and href paths.
func ListenAndServe(dir string) {
	log.Println("listening to 127.0.0.1:8080")
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func Must(ws *Website, err error) *Website {
	if err != nil {
		panic(err)
	}
	return ws
}
