package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	// Handle environment config values
	addr := flag.String("addr", ":4000", "HTTP port")
	flag.Parse()

	// Initialise loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	// Set up application struct
	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

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

	// Create HTTP server struct
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	infoLog.Printf("Starting server on http://localhost:%s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
