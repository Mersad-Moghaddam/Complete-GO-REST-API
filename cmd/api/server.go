package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// server starts a new HTTP server listening on the port specified in the
// application configuration.
func (app *Application) serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on port %d", app.port)

	return server.ListenAndServe()
}
