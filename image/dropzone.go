package image

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"hash/adler32"
	"html/template"
	"net/http"
	"path"
)

//go:embed dropzone/*
var dropzoneStatic embed.FS // contains "dropzone" dir

//go:embed *.html
var htmlfiles embed.FS
var tmpl = template.Must(template.ParseFS(htmlfiles, "dropzone.html"))

type tmplData struct {
	APIPath      string
	ID           string
	Images       []Image
	DeleteURL    string
	GetURLPrefix string
	UploadURL    string
}

// DropzoneHandler handles API calls from DropzoneHTML and serves dropzone CSS and JS. You have to take care for authentication.
func (store Store) DropzoneHandler(onModify func()) http.Handler {
	var mux = http.NewServeMux()
	mux.HandleFunc("POST /delete/{path...}", func(w http.ResponseWriter, r *http.Request) {
		filename := r.PostFormValue("name")
		if err := store.Delete(r.PathValue("path"), filename); err != nil {
			http.Error(w, fmt.Sprintf("error removing image: %v", err), http.StatusInternalServerError)
			return
		}
		if onModify != nil {
			onModify()
		}
	})
	mux.Handle("GET /get/{path...}", http.StripPrefix("/get", Store{ // copy Store so we can set MaxSides
		CacheDir: store.CacheDir,
		Dir:      store.Dir,
		MaxSides: []int{300}, // 300 is used in dropzone.html
	}))
	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.FS(dropzoneStatic))))
	mux.HandleFunc("POST /upload/{path...}", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file") // name="file" from dropzonejs
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting form file: %v", err), http.StatusInternalServerError)
			return
		}
		if err := store.Upload(r.PathValue("path"), file, header); err != nil {
			http.Error(w, fmt.Sprintf("error storing image: %v", err), http.StatusInternalServerError)
			return
		}
		if onModify != nil {
			onModify()
		}
	})
	return mux
}

func (store Store) DropzoneHTML(apiPath, dir string) (template.HTML, error) {
	hash := adler32.New()
	hash.Write([]byte(apiPath + "|" + dir))
	id := base64.RawURLEncoding.EncodeToString(hash.Sum(nil)) // RawURLEncoding (dash, underscore, no padding) is okay for HTML id value

	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, tmplData{
		APIPath:      path.Join("/", apiPath, "/"),
		ID:           id,
		Images:       store.Get(dir),
		DeleteURL:    path.Join("/", apiPath, "delete", dir),
		GetURLPrefix: path.Join("/", apiPath, "get", dir),
		UploadURL:    path.Join("/", apiPath, "upload", dir),
	})
	return template.HTML(buf.String()), err
}
