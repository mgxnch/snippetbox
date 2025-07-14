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
	// Since serverError is a helper (call depth 1), we don't want this to be the source
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
