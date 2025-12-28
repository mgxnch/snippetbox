package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/mgxnch/snippetbox/internal/models"
	"github.com/mgxnch/snippetbox/ui"
)

// templateData is a holding struct for data that needs to be passed to
// ExecuteTemplate(), allowing us to pass multiple data fields into
// ExecuteTemplate, which only accepts a single data object in its parameters.
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any    // holds validation errors
	Flash           string // holds the flash message
	IsAuthenticated bool   // true if user is authenticated, false otherwise
	CSRFToken       string // holds the CSRF token
}

// functions acts as a lookup between the names of our custom template
// functions and the functions. Custom template functions must only
// return one value, or two values where the second value is an error.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Init the map
	cache := map[string]*template.Template{}

	// fs.Glob returns a slice of filepath strings that match the pattern
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (e.g. 'home.tmpl') from the full filepath
		name := filepath.Base(page)

		// Define a slice of filepath patterns for the templates that
		// we want to parse. Each of the .tmpl files in ui/html/pages require
		// base.tmpl and all of the partials' .tmpl files.
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// ParseFS is a variadic function, which allows us to parse multiple templates in a single
		// call. We no longer have to split between ParseFiles and ParseGlob.
		tmpl, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Use the base file name as map key for the parsed template
		cache[name] = tmpl
	}

	return cache, nil
}

// humanDate is used as a template function which formats a time.Time
// struct into a human-readable string.
func humanDate(t time.Time) string {
	// This time is Go's reference time
	// ref: https://stackoverflow.com/questions/28087471/what-is-the-significance-of-gos-time-formatlayout-string-reference-time
	return t.Format("02 Jan 2006 15:04")
}
