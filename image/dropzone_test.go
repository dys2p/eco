package image

import "net/http"

func Example() {
	store := &Store{
		CacheDir: "/tmp/cache",
		Dir:      "/tmp/images",
		MaxSides: []int{1200},
	}
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		h, err := store.DropzoneHTML("/", "/")
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html>" + h + "</html>"))
	})
	http.Handle("/", store.DropzoneHandler(nil))
	http.ListenAndServe("127.0.0.1:8080", nil)
}
