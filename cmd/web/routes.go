package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mgxnch/snippetbox/ui"
)

// routes sets up a Chi router, its routes and returns the object.
func (app *application) routes() http.Handler {
	// Initialise Chi router
	r := chi.NewRouter()

	// Add custom 404 handler
	// ref: https://go-chi.io/#/pages/routing?id=making-custom-404-and-405-handlers
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Middleware chain:
	// recoverPanic ->logRequest -> secureHeaders -> serverMux -> application handlers
	// Chi middlewares have to be declared before routes
	r.Use(app.recoverPanic)
	r.Use(app.logRequest)
	r.Use(secureHeaders)

	// File server and its route
	// Convert ui.Files from an embedded filesystem to the http.FS type, so that
	// it satisfies the http.FileSystem type required of the argument to http.FileServer.
	fileServer := http.FileServer(http.FS(ui.Files))

	// ui.Files has two subdirectories - "html" and "static". Any requests with the "/static"
	// prefix now exists within the file server.
	r.Handle("/static/*", fileServer)

	// Application routes that use the Session Manager
	// We use r.Group if not Chi will complain that we are declaring middleware
	// components after routes. Chi only allows you to declare middleware BEFORE routes.
	r.Group(func(r chi.Router) {
		// Add the middleware for this group
		r.Use(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

		// Add the handlers for this group
		r.Get("/", app.home)
		r.Get("/snippet/view/{id}", app.snippetView)
		r.Get("/user/signup", app.userSignup)
		r.Post("/user/signup", app.userSignupPost)
		r.Get("/user/login", app.userLogin)
		r.Post("/user/login", app.userLoginPost)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthentication)

			r.Get("/snippet/create", app.snippetCreate)
			r.Post("/snippet/create", app.snippetCreatePost)
			r.Post("/user/logout", app.userLogoutPost)
		})
	})

	return r
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deferred functions will run right before the actual return is completed
		defer func() {
			// Check if there has been a panic
			if err := recover(); err != nil {
				// Close the connection and return a server error
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
