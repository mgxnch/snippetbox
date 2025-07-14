package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Handle environment config values
	addr := flag.String("addr", ":4000", "HTTP port")
	flag.Parse()

	// Initialise router
	mux := http.NewServeMux()

	// File server and its route
	// http.FileServer serves file out of its input directory
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Strip away /static from the incoming request e.g. GET localhost:4000/static/img/logo.png
	// becomes "img/logo.png" which is a valid path to the fileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Application routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on http://localhost:%s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
