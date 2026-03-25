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
// You can return an http.RedirectHandler. You can also return an http.NotFoundHandler, although it's probably better to write your own error handlers, e. g.:
//
//	func internalServerError(err error) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			// log error
//			w.WriteHeader(http.StatusInternalServerError)
//			// execute custom error html template
//		})
//	}
//
//	func executeTemplate(t *template.Template, data any) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			if err := t.Execute(w, data); err != nil {
//				internalServerError(err).ServeHTTP(w, r)
//			}
//		})
//	}
type HandlerFunc func(w http.ResponseWriter, r *http.Request) http.Handler

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := f(w, r); handler != nil {
		handler.ServeHTTP(w, r)
	}
}
