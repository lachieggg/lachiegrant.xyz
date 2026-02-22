package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareFunc(t *testing.T) {
	// Setup a clean logger specifically bounded to a buffer so we can parse logs.
	var logOutput bytes.Buffer
	logger = log.New(&logOutput, "", 0)

	// Create a generic handler representing any of our downstream functions (blog, resume, index)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot) // Return 418 so we know it specifically triggered inside here.
		w.Write([]byte("Teapot handler triggered"))
	})

	// Wrap our basic handler exactly as done in main()
	wrappedHandler := middlewareFunc(testHandler)

	tests := []struct {
		name        string
		method      string
		url         string
		remoteIP    string
		expectedLog string
	}{
		{
			name:        "Formats valid generic GET request from Localhost",
			method:      http.MethodGet,
			url:         "http://localhost/test-path",
			remoteIP:    "127.0.0.1",
			expectedLog: "Request: GET /test-path from 127.0.0.1 - Status: 418",
		},
		{
			name:        "Formats valid POST request from External Entity",
			method:      http.MethodPost,
			url:         "http://localhost/data-push",
			remoteIP:    "192.168.1.1",
			expectedLog: "Request: POST /data-push from 192.168.1.1 - Status: 418",
		},
		{
			name:        "Gracefully handles missing X-Real-IP mapping naturally",
			method:      http.MethodOptions,
			url:         "http://localhost/empty-ip",
			remoteIP:    "",
			expectedLog: "Request: OPTIONS /empty-ip from  - Status: 418",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logOutput.Reset()

			req := httptest.NewRequest(tc.method, tc.url, nil)
			if tc.remoteIP != "" {
				req.Header.Set("X-Real-IP", tc.remoteIP)
			}

			rec := httptest.NewRecorder()

			wrappedHandler(rec, req)

			// 1. Assert inner handler successfully ran
			assert.Equal(t, http.StatusTeapot, rec.Code, "Response Code matching Teapot expected")
			assert.Equal(t, "Teapot handler triggered", rec.Body.String(), "Body matching Teapot expected")

			// 2. Assert log structurally represents exactly what the middleware was expected to see
			assert.Contains(t, logOutput.String(), tc.expectedLog, "Middleware Log Output should exactly map the input criteria")
		})
	}
}
