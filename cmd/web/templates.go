package main

import "github.com/mgxnch/snippetbox/internal/models"

// templateData is a holding struct for data that needs to be passed to
// ExecuteTemplate(), allowing us to pass multiple data fields into
// ExecuteTemplate, which only accepts a single data object in its parameters.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
