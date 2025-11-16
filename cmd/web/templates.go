package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/mgxnch/snippetbox/internal/models"
)

// templateData is a holding struct for data that needs to be passed to
// ExecuteTemplate(), allowing us to pass multiple data fields into
// ExecuteTemplate, which only accepts a single data object in its parameters.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any    // holds validation errors
	Flash       string // holds the flash message
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

	// filepath.Glob returns a slice of filepath strings that match the pattern
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (e.g. 'home.tmpl') from the full filepath
		name := filepath.Base(page)

		// Instead of passing the full list of files to teplate.ParseFiles,
		// we can create the Template object first, register the template
		// functions that we want, then start adding more files to it
		tmpl := template.New(name).Funcs(functions)

		// Parse base.tmpl into a template object
		tmpl, err := tmpl.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Add partials to the template
		// note(mx): we are calling Parse[Glob|Files] on our Template object
		tmpl, err = tmpl.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Add the current page into the template
		tmpl, err = tmpl.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// and use it as the map key for the parsed template
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
