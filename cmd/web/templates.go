package main

import (
	"html/template"
	"path/filepath"

	"github.com/mgxnch/snippetbox/internal/models"
)

// templateData is a holding struct for data that needs to be passed to
// ExecuteTemplate(), allowing us to pass multiple data fields into
// ExecuteTemplate, which only accepts a single data object in its parameters.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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
		// Instead of passing the full list of files to teplate.ParseFiles,
		// we can create the Template object first, then start adding more
		// files to it

		// Parse base.tmpl into a template object
		tmpl, err := template.ParseFiles("./ui/html/base.tmpl")
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

		// Extract the file name (e.g. 'home.tmpl') from the full filepath
		// and use it as the map key for the parsed template
		name := filepath.Base(page)
		cache[name] = tmpl
	}

	return cache, nil
}
