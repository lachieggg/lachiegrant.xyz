package main

import (
	"fmt"
	"net/http"
	"net/url"
)

// middlewareFunc is a general middleware function
// for all endpoints
func middlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Log the request
		if logger == nil || r == nil {
			fmt.Println("Cannot log request")
			fmt.Printf("logger: %v\n", logger)
			fmt.Printf("request: %v\n", r)
			next(w, r)
			return
		}

		remoteAddr := r.Header.Get("X-Real-IP")

		if r.URL == nil {
			logger.Printf("URL is nil, remote address: %s\n", remoteAddr)
			r.URL = &url.URL{}
		}

		logger.Printf("Request received for endpoint %s from %s\n",
			r.URL.Path,
			remoteAddr,
		)

		next(w, r)
	}
}
