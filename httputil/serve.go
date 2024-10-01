package httputil

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

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
func ListenAndServe(addr string, handler http.Handler, stop chan<- os.Signal) func() {
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
