package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
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

// render is a helper to render a HTML page and return it to the client.
// page is the base file path of the *.tmpl files in the "ui/html/pages/" folder
// e.g. "home.tmpl"
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
