package main

import "net/http"

// middlewareFunc is a general middleware function
// for all endpoints
func middlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        // Log the request
		if logger != nil && r != nil && r.URL != nil {
			logger.Printf("Request received for endpoint %s from %s", r.URL.Path, r.RemoteAddr)
		}
        next(w, r)
    }
}
