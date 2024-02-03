// Package httputil provides an easy way to chain handlers and a server with timeouts and graceful shutdown.
package httputil

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
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

// ListenAndServe listens on the TCP network address addr and starts an HTTP server with timeouts.
// It returns a shutdown function for the server.
//
// Example:
//
//	stop := make(chan os.Signal, 1)
//	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
//
//	shutdown := httputil.ListenAndServe(":8080", router, stop)
//	defer shutdown()
//
//	<-stop
func ListenAndServe(addr string, handler http.Handler, stop chan os.Signal) func() {
	// "You should set Read, Write and Idle timeouts when dealing with untrusted clients and/or networks"
	// https://blog.cloudflare.com/exposing-go-on-the-internet/
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  30 * time.Second, // client sending request
		WriteTimeout: 30 * time.Second, // server reading request body + writing response
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		// http.Server.ListenAndServe creates a tcpKeepAliveListener
		if err := srv.ListenAndServe(); err != http.ErrServerClosed { // ErrServerClosed is ok
			log.Printf("error listening: %v", err)
			stop <- os.Interrupt // SIGINT
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("error shutting down: %v", err)
		}
	}
}
