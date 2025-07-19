package main

import (
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

// render is a helper to render a HTML page and return it to the client. The page
// argument is the base file path of the *.tmpl files in the "ui/html/pages/" folder
// e.g. "home.tmpl"
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	tmpl, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Useful for when we need to customise the status. Typically we can
	// set this to 200 OK by default
	w.WriteHeader(status)

	// ExecuteTemplate writes the output of the parsed template into the writer w
	// We have to specify the named template to parse and apply
	// Template names are declared in the .tmpl file in the {{define "xxx"}} block
	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}
