package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // import for side-effects only
	"github.com/mgxnch/snippetbox/internal/models"
)

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder // for validating form fields
	sessionManager *scs.SessionManager
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

	// Set up session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Set up application struct
	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{
			DB: db,
		},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Set up non-default TLS settings. We are using these two with assembly implementations
	// which should in theory be much faster. tls.Config supports other preferences too,
	// such as CipherSuites, but we are not setting that for now.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Create HTTP server struct
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
	}
	infoLog.Printf("Starting server on http://localhost%s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
