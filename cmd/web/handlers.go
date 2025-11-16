package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mgxnch/snippetbox/internal/models"
	"github.com/mgxnch/snippetbox/internal/validator"
)

// home is the function handler for the root page. It fetches the latest 10
// snippets and renders to the user.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
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
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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

	// Reads and removes the 'flash' key from sessionManager
	flash := app.sessionManager.PopString(r.Context(), "flash")

	// Populate the templateData struct with data
	data := newTemplateData()
	data.Snippet = snippet
	data.Flash = flash

	// Render the page
	app.render(w, http.StatusOK, "view.tmpl", data)
}

// snippetCreateForm represents the form data and validation errors
// for the snippetCreate form fields.
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // embedded struct
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := newTemplateData()
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	// Read submitted form fields
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validation logic for title, content and expires
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// If there are any validation errors, re-render the create.tmpl template
	if !form.Valid() {
		data := newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a key and its string value to the session data
	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully")

	// If snippet is successfully created, redirect user to the snippetView page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
