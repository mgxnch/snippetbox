package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
)

// routes sets up a ServeMux and its routes and returns the object.
func (app *application) routes() http.Handler {
	// Initialise router
	mux := http.NewServeMux()

	// File server and its route
	// http.FileServer serves file out of its input directory
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Strip away /static from the incoming request e.g. GET localhost:4000/static/img/logo.png
	// becomes "img/logo.png" which is a valid path to the fileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Application routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Middleware chain: recoverPanic ->logRequest -> secureHeaders -> serverMux -> application handlers

	// Create a 'standard' middleware chain which contains the middlewares
	// that will be used for every request our application receives
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
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
