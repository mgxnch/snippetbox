package main

import (
	"context"
	"net/http"

	"github.com/justinas/nosurf"
)

// secureHeaders is a middleware that sets security-related headers
// into the HTTP response in accordance with OWASP best practices.
func secureHeaders(next http.Handler) http.Handler {
	// Standard middleware convention
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src: 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Xss-Protection", "0")

		// Important thing is to call next.ServeHTTP() to keep the chain going
		next.ServeHTTP(w, r)

		// If you want code to execute on the way back up the chain,
		// you include them after the next.ServeHTTP call
	})
}

// noSurf is a middleware which uses a customised CSRF cookie with the
// Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

// requireAuthentication is a middleware that checks if a user is allowed to
// access a certain page.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// So that pages that require authentication are not stored in the user's
		// browser cache, or other intermediary caches
		w.Header().Add("Cache-Control", "no-store")

		// Call next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// authenticate is a middleware that checks if a userID exists within the context.
// If a valid userID exists, the context is set with the isAuthenticatedContextKey
// key with a value of true.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), authUserKey)
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Check DB to see if userID exists
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
