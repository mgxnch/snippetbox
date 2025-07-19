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
		// Declare files to be parsed by template.ParseFiles
		// Create a slice containing the filepaths for our base template,
		// any partials and the page itself. Base and partials are always
		// included in every page, hence we always include their filepaths
		// for every page.
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		// Parse the files into a template
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// Extract the file name (e.g. 'home.tmpl') from the full filepath
		name := filepath.Base(page)

		// Set the key as name and value as the template
		cache[name] = ts
	}

	return cache, nil
}
