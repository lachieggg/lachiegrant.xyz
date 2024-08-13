package main

import (
	"net/http"
	"net/url"
)

// middlewareFunc is a general middleware function
// for all endpoints
func middlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.Header.Get("X-Real-IP")

		if r.URL == nil {
			r.URL = &url.URL{}
		}

		logger.Printf("Request received for endpoint %s from %s\n",
			r.URL.Path,
			remoteAddr,
		)

		next(w, r)
	}
}
