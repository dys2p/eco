// Package middleware provides a net/http compatible way to chain http handlers.
//
//	mux.Handle("/", middleware.HandlerFunc(func(w http.ResponseWriter, r *http.Request) http.Handler {
//		// do something, then return next handler or nil
//	}))
//
// Use http.NotFound and http.Redirect for not found errors and redirects.
package middleware

import (
	"log"
	"net/http"
)

// A HandlerFunc can return the next handler or nil.
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
