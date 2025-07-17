package main

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/mgxnch/snippetbox/internal/models"
)

// home is a function haldner which writes a byte slice
// as the response body
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Prevent non-existent paths from matching on this handler
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	// ts stands for Template Set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// ExecuteTemplate writes the output of the parsed template into the writer w
	// We have to specify the named template to parse and apply
	// Template names are declared in the .tmpl file in the {{define "xxx"}} block
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(("id")))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Fetch the snippet by its ID
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Set up template
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl", // view.tmpl contains a "main" named template
	}

	// Parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Pass a pointer to the templateData as the data argument to ExecuteTemplate
	err = ts.ExecuteTemplate(w, "base", &templateData{Snippet: snippet})
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Allow
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("create a new snippet"))
}
