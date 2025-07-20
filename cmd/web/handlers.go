package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/mgxnch/snippetbox/internal/models"
)

// home is the function handler for the root page. It fetches the latest 10
// snippets and renders to the user.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Prevent non-existent paths from matching on this handler
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Fetch all snippets from DB
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	// Populate the templateData struct with data
	data := newTemplateData()
	data.Snippets = snippets

	// Render the page
	app.render(w, http.StatusOK, "home.tmpl", data)
}

// snippetView is the function handler for viewing a specific snippet.
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

	// Populate the templateData struct with data
	data := newTemplateData()
	data.Snippet = snippet

	// Render the page
	app.render(w, http.StatusOK, "view.tmpl", data)
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
