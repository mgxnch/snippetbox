package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // import for side-effects only
	"github.com/mgxnch/snippetbox/internal/models"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder // for validating form fields
}

func main() {
	// Handle environment config values
	addr := flag.String("addr", ":4000", "HTTP port")
	dsn := flag.String("dsn", "web:9mfOz8RWTWQSIlgt8hX9jb9V@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Initialise loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialise DB pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Set up the template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Set up a decoder instance
	formDecoder := form.NewDecoder()

	// Set up application struct
	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{
			DB: db,
		},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// Create HTTP server struct
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on http://localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check that database is connected
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
