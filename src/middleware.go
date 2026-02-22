package main

import "net/http"

// responseWriter is a wrapper around http.ResponseWriter that captures the
// HTTP status code of the response.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// middlewareFunc wraps an HTTP handler with request logging.
// It logs the requested path, the client's IP address, and the response status code.
func middlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// X-Real-IP is set by the reverse proxy to preserve the original client IP
		remoteAddr := r.Header.Get("X-Real-IP")

		next(rw, r)

		logger.Printf("Request: %s %s from %s - Status: %d", r.Method, r.URL.Path, remoteAddr, rw.status)
	}
}
