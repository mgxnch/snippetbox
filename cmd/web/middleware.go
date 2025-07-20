package main

import "net/http"

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
