package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// routes sets up a Chi router, its routes and returns the object.
func (app *application) routes() http.Handler {
	// Initialise Chi router
	r := chi.NewRouter()

	// Add custom 404 handler
	// ref: https://go-chi.io/#/pages/routing?id=making-custom-404-and-405-handlers
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Middleware chain:
	// recoverPanic ->logRequest -> secureHeaders -> serverMux -> application handlers
	// Chi middlewares have to be declared before routes
	r.Use(app.recoverPanic)
	r.Use(app.logRequest)
	r.Use(secureHeaders)

	// File server and its route
	// http.FileServer serves file out of its input directory
	// "." refers the current working directory, i.e. the directory
	// where you're running the binary from, and NOT where the Go binary resides on disk
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Strip away /static from the incoming request e.g. GET localhost:4000/static/img/logo.png
	// becomes "img/logo.png" which is a valid path to the fileServer
	r.Handle("/static/*", http.StripPrefix("/static", fileServer)) // http.StripPrefix returns a handler

	// Application routes
	r.Get("/", app.home)
	r.Get("/snippet/view/{id}", app.snippetView)
	r.Get("/snippet/create", app.snippetCreate)
	r.Post("/snippet/create", app.snippetCreatePost)

	return r
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deferred functions will run right before the actual return is completed
		defer func() {
			// Check if there has been a panic
			if err := recover(); err != nil {
				// Close the connection and return a server error
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
