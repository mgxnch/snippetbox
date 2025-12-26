package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

// serverError is a helper to print the error stack trace and return HTTP 500 to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	// We set calldepth to 2, which represents caller of the caller
	// Since serverError is a helper (call depth 1), we don't want this call site
	// to be the source when we are looking at the stack trace
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// notFound is a helper to return 404 Not Found to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// clientError is a helper to return client-related HTTP errors e.g. 400 Bad Request
func (app *application) clientError(w http.ResponseWriter, status int) {
	// note(mx): http.Error calls w.Write downstream
	http.Error(w, http.StatusText(status), status)
}

// render writes a HTML template with a name of page to w.
//
// page is the base file path of the *.tmpl files in the "ui/html/pages/" folder
// e.g. "home.tmpl"
//
// Note: render only writes to w and is not responsible for returning a page to the user.
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	tmpl, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Intermediate buffer to store the result of ExecuteTemplate
	// instead of immediately writing to http.ResponseWriter
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Useful for when we need to customise the status. Typically we can
	// set this to 200 OK by default
	w.WriteHeader(status)
	buf.WriteTo(w)
}

// decodePostForm is a helper to to extract the request's form fields by using
// a form decoder. If the form is well-formed, its values are set into dst. dst
// is expected to be a pointer to a struct that can hold the form values.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// r.ParseForm populates r.Form and r.PostForm. r.PostForm can be seen as a subset
	// of r.Form, as it only contains values from POST, PUT and PATCH requests
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Use Decode to parse r.PostForm and set it into our form struct
	// Type conversions are automatically handled for us
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// InvalidDecoderError is raised when a nil pointer is passed to Decode
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal
		return err
	}
	return nil
}
