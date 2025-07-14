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
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Set up application struct
	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	// Create HTTP server struct
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on http://localhost:%s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
