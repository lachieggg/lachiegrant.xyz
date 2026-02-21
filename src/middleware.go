package main

import "net/http"

// middlewareFunc wraps an HTTP handler with request logging.
// It logs the requested path and the client's IP address (from X-Real-IP header,
// typically set by a reverse proxy like nginx).
func middlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// X-Real-IP is set by the reverse proxy to preserve the original client IP
		remoteAddr := r.Header.Get("X-Real-IP")

		logger.Printf("Request: %s from %s", r.URL.Path, remoteAddr)

		next(w, r)
	}
}
