// Package image provides an image store and resizer. It depends on convert (from imagemagick). Uploaded files are stored in original, but not served via HTTP.
//
//	store := image.Store{
//	    CacheDir: filepath.Join(os.Getenv("CACHE_DIRECTORY"), "imagedir"),
//	    Dir:      filepath.Join(os.Getenv("STATE_DIRECTORY"), "imagedir"),
//	    MaxSides: []int{300, 600, 1200},
//	}
//
//	router.Handler(http.MethodGet, "/images/*filepath", http.StripPrefix("/images", store))
package image

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type Store struct {
	CacheDir    string
	Dir         string
	JPEGQuality int // default: 80
	MaxSides    []int
}

func (s Store) jpegQuality() string {
	if s.JPEGQuality >= 50 && s.JPEGQuality <= 100 {
		return strconv.Itoa(s.JPEGQuality)
	} else {
		return "80"
	}
}

func (store Store) Get(dir string) []Image {
	dir = filepath.Join("/", dir) // removes "..", preventing directory traversal
	var result []Image
	entries, _ := os.ReadDir(filepath.Join(store.Dir, dir))
	for _, entry := range entries {
		image := Image{
			Store: store,
			Dir:   dir,
			Name:  entry.Name(),
		}
		if info, err := entry.Info(); err == nil {
			image.Size = info.Size()
		}
		result = append(result, image)
	}
	return result
}

func (store Store) Remove(dir, filename string) error {
	dir = filepath.Join("/", dir) // removes "..", preventing directory traversal
	if strings.ContainsRune(filename, filepath.Separator) {
		return errors.New("invalid filename")
	}

	// remove original
	if err := os.Remove(filepath.Join(store.Dir, dir, filename)); err != nil {
		return err
	}
	// remove cached
	return os.RemoveAll(filepath.Join(store.CacheDir, dir, filename))
}

func (store Store) Upload(dir string, file multipart.File, header *multipart.FileHeader) error {
	defer file.Close()

	if header.Size > 20*1024*1024 {
		return errors.New("file too large")
	}

	dir = filepath.Join("/", dir) // removes "..", preventing directory traversal
	if strings.ContainsRune(header.Filename, filepath.Separator) {
		return errors.New("filename contains path separator")
	}

	fp := filepath.Join(store.Dir, dir, header.Filename)
	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return fmt.Errorf("making directory: %w", err)
	}

	// copy file
	osFile, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("creating or opening target file for writing: %w", err)
	}
	if _, err := io.Copy(osFile, file); err != nil {
		osFile.Close()
		return fmt.Errorf("writing to target file: %w", err)
	}
	osFile.Close()

	return nil
}

func (store Store) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	maxSide, _ := strconv.Atoi(r.URL.Query().Get("m"))
	if maxSide <= 0 || !slices.Contains(store.MaxSides, maxSide) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// filepaths
	reqpath := filepath.FromSlash(path.Join("/", r.URL.Path)) // "/" ensures the path is absolute, preventing directory traversal
	orig := filepath.Join(store.Dir, reqpath)
	cached := filepath.Join(store.CacheDir, orig, "max-side", strconv.Itoa(maxSide)) // path element order makes cache invalidation easy
	if filepath.Ext(cached) != ".jpg" {
		cached = cached + ".jpg" // imagemagick determines file format from extension
	}

	// check if original exists
	if _, err := os.Lstat(orig); err == nil {
		// continue
	} else if errors.Is(err, os.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if cached exists
	if _, err := os.Lstat(cached); err == nil {
		http.ServeFile(w, r, cached)
		return
	} else if errors.Is(err, os.ErrNotExist) {
		// continue
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create cached
	if err := os.MkdirAll(filepath.Dir(cached), 0755); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := exec.Command("convert", "-resize", fmt.Sprintf("%dx%d>", maxSide, maxSide), "-quality", store.jpegQuality(), "-alpha", "remove", "-background", "white", orig, cached).Run(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, cached)
}
