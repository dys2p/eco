// Package httputil provides an easy way to chain handlers and a server with timeouts and graceful shutdown.
package httputil

import (
	"log"
	"net/http"
)

// A HandlerFunc can return the next handler or nil, which provides easy chaining:
//
//	mux.Handle("/", httputil.HandlerFunc(func(w http.ResponseWriter, r *http.Request) http.Handler {
//		// do something, then return next handler or nil
//	}))
//
// You can return http.NotFoundHandler and http.RedirectHandler for not found errors and redirects.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) http.Handler

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := f(w, r); handler != nil {
		handler.ServeHTTP(w, r)
	}
}

func Forbidden() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
}

// InternalServerError returns an handler which writes status code 500 and logs the error string.
func InternalServerError(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("internal server error: %v", err)
	})
}
