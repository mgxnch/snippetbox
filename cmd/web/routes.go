package main

import "net/http"

// routes sets up a ServeMux and its routes and returns the object.
func (app *application) routes() *http.ServeMux {
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

	return mux
}
