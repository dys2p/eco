// Package ssg serves multi-lingual websites. It's not really a static site generator because HTTP redirects can't be expressed as HTML files.
// Content is stored in markdown files. Symlinks are supported. Use with trusted input only.
package ssg

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/dys2p/eco/lang"
	"gitlab.com/golang-commonmark/markdown"
	"golang.org/x/net/html"
	"golang.org/x/text/language"
)

var md = markdown.New(markdown.HTML(true), markdown.Linkify(false))

// LangOption should be used in templates.
type LangOption struct {
	BCP47    string
	Name     string
	Prefix   string
	Selected bool
}

// LangOptions returns an empty slice if there is only one language.
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

// TemplateData is recommended for templates. Custom template data can embed it.
type TemplateData struct {
	lang.Lang
	Languages []LangOption
	Path      string // without language prefix, for language buttons and hreflang
	Title     string // for <title>
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

// one markdown file
type translation struct {
	tag     language.Tag
	content string // HTML
}

// Make returns a HTTP handler which serves content from fsys.
// It optionally accepts an additional HTML template and a function which makes custom template data.
func Make(fsys fs.FS, add *template.Template, langs lang.Languages, makePageData func(*http.Request, TemplateData) any, notFound http.Handler) (http.Handler, error) {
	var mux = http.NewServeMux()
	// walk fsys
	var langPrefixes = make(map[string]any)
	for _, l := range langs {
		langPrefixes[l.Prefix] = struct{}{}
	}
	var pageNames = make(map[string]any)
	err := walk(fsys, ".", mux, langPrefixes, pageNames)
	if err != nil {
		return nil, fmt.Errorf("walking filesystem: %w", err)
	}
	// parse and assemble common template
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
	// makePageData default value
	if makePageData == nil {
		makePageData = func(_ *http.Request, data TemplateData) any {
			return data
		}
	}
	// translate pages
	for pageName := range pageNames {
		// read translated markdown files, if exist
		var translations []translation
		for _, l := range langs {
			bs, err := fs.ReadFile(fsys, filepath.Join(pageName, l.Prefix+".md"))
			if err == nil {
				translations = append(translations, translation{
					tag:     l.Tag,
					content: md.RenderToString(bs),
				})
			}
		}
		if len(translations) == 0 {
			continue
		}
		// make matcher for available translations
		var translationTags []language.Tag
		for _, tr := range translations {
			translationTags = append(translationTags, tr.tag)
		}
		matcher := language.NewMatcher(translationTags)
		// handle redirect to localized url (using Accept-Language header)
		mux.HandleFunc("GET /"+pageName+".html", func(w http.ResponseWriter, r *http.Request) {
			_, index := language.MatchStrings(matcher, r.Header.Get("Accept-Language"))
			var u = *r.URL // copy
			u.Path = path.Join("/", langs[index].Prefix, u.Path)
			http.Redirect(w, r, u.String(), http.StatusSeeOther)
		})
		// handle page for all languages
		for _, l := range langs {
			// clone common template
			tt, err := tmpl.Clone()
			if err != nil {
				return nil, fmt.Errorf("cloning template: %w", err)
			}
			// insert translated content
			_, index, _ := matcher.Match(l.Tag)
			tt, err = tt.Parse(`{{define "content"}}` + translations[index].content + `{{end}}`) // or parse content into t and then call AddParseTree(content, t.Tree)
			if err != nil {
				return nil, fmt.Errorf("adding content of %s: %w", pageName, err)
			}
			// handler calls makePageData and executes template
			pattern := "GET " + path.Join("/", l.Prefix, pageName+".html")
			mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
				data := makePageData(r, TemplateData{
					Lang:      l,
					Languages: LangOptions(langs, l),
					Path:      pageName + ".html",
					Title:     heading(translations[index].content),
				})
				if err := tt.ExecuteTemplate(w, "html", data); err != nil {
					log.Printf("executing html template %s: %v", pattern, err)
				}
			})
		}
	}
	// chain next handler or redirect to index.html if exists
	if notFound != nil {
		mux.Handle("/", notFound)
	} else if _, ok := pageNames["index"]; ok {
		mux.Handle("/", http.RedirectHandler("/index.html", http.StatusSeeOther))
	}
	return mux, nil
}

// concatenates all text nodes within the first h1/h2/h3/h4 node
func heading(htm string) string {
	r := &io.LimitedReader{
		R: strings.NewReader(htm),
		N: 4096,
	}
	var tokenizer = html.NewTokenizerFragment(r, "body")
	var headingTag = ""
	var result strings.Builder
nextToken:
	for {
		tt := tokenizer.Next()
		tnb, _ := tokenizer.TagName()
		tn := string(tnb)
		switch {
		case tt == html.ErrorToken:
			break nextToken // EOF
		case headingTag == "" && tt == html.StartTagToken && (tn == "h1" || tn == "h2" || tn == "h3" || tn == "h4"):
			headingTag = tn // <h1>
		case headingTag != "" && tt == html.EndTagToken && tn == headingTag:
			break nextToken // </h1>
		case headingTag != "":
			result.Write(tokenizer.Raw()) // between <h1> and </h1>
		}
	}
	return result.String()
}

func walk(fsys fs.FS, dir string, mux *http.ServeMux, langPrefixes map[string]any, pageNames map[string]any) error {
	dir = strings.Trim(dir, "/") // io/fs docs: dir must not start or end with a slash
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		entryPath := path.Join(dir, entry.Name())
		// skip dot entries
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		// check if symlink is dir
		var isDir = entry.IsDir()
		if entry.Type()&fs.ModeSymlink != 0 {
			if info, err := fs.Stat(fsys, entryPath); err == nil { // fs.Stat follows symlink
				if info.Mode()&fs.ModeDir != 0 {
					isDir = true
				}
			}
		}
		// recursion
		if isDir {
			if err := walk(fsys, entryPath, mux, langPrefixes, pageNames); err != nil {
				return err
			}
			continue
		}
		// skip html files in root dir (process them later)
		ext := path.Ext(entry.Name())
		if (dir == "" || dir == ".") && ext == ".html" {
			continue
		}
		// lang.md establishes a page
		root := strings.TrimSuffix(entry.Name(), ext)
		_, rootIsLangPrefix := langPrefixes[root]
		if rootIsLangPrefix && ext == ".md" {
			pageNames[dir] = struct{}{}
			continue
		}
		// serve everything else as static file
		mux.HandleFunc("GET "+path.Join("/", entryPath), func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, fsys, path.Join("/", entryPath))
		})
	}
	return nil
}
