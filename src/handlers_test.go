package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMain sets up the test environment, including initializing the global logger.
func TestMain(m *testing.M) {
	// Initialize logger to discard output during tests
	logger = log.New(io.Discard, "", 0)

	// go test runs natively inside the folder containing the tests (./src).
	// We dynamically step back one directory so templates route logically to ./src/templates.
	os.Chdir("..")

	os.Exit(m.Run())
}

func TestHandlers(t *testing.T) {
	// Setup: Create a dummy public directory and file for resume tests
	err := os.MkdirAll("public", 0755)
	assert.NoError(t, err)

	testResumeFile := "public/test-resume.pdf"
	testContent := "ok"
	err = os.WriteFile(testResumeFile, []byte(testContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testResumeFile)

	tests := []struct {
		name           string
		method         string
		url            string
		handler        http.HandlerFunc
		envVars        map[string]string
		expectedStatus int
		expectedBody   string
		expectedHeader map[string]string
	}{
		// Index Handler Tests
		{
			name:           "Page: Home (Exists)",
			url:            "/",
			handler:        indexHandler,
			expectedStatus: http.StatusOK, // Note: template-dependent, usually 200/500 in test
		},
		{
			name:           "Page: Not Found (Simple)",
			url:            "/nonexistent",
			handler:        indexHandler,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Page: Not Found (Nested)",
			url:            "/nested/path",
			handler:        indexHandler,
			expectedStatus: http.StatusNotFound,
		},

		// GitHub Handler Tests
		{
			name:           "Redirect: GitHub (Valid)",
			url:            "/code",
			handler:        githubHandler,
			envVars:        map[string]string{EnvGithubURL: "https://github.com/lachieggg"},
			expectedStatus: http.StatusSeeOther,
			expectedHeader: map[string]string{"Location": "https://github.com/lachieggg"},
		},
		{
			name:           "Redirect: GitHub (Empty)",
			url:            "/code",
			handler:        githubHandler,
			envVars:        map[string]string{EnvGithubURL: ""},
			expectedStatus: http.StatusSeeOther,
		},

		// Resume Handler Tests
		{
			name:           "File: Resume (Success)",
			url:            "/resume",
			handler:        resumeHandler,
			envVars:        map[string]string{EnvResumePath: "test-resume.pdf"},
			expectedStatus: http.StatusOK,
			expectedBody:   testContent,
		},
		{
			name:           "File: Resume (Traversal Attack)",
			url:            "/resume",
			handler:        resumeHandler,
			envVars:        map[string]string{EnvResumePath: "../main.go"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "File: Resume (System Path Attack)",
			url:            "/resume",
			handler:        resumeHandler,
			envVars:        map[string]string{EnvResumePath: "/etc/passwd"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "File: Resume (Not Found)",
			url:            "/resume",
			handler:        resumeHandler,
			envVars:        map[string]string{EnvResumePath: "missing.pdf"},
			expectedStatus: http.StatusNotFound,
		},

		// Feature Flag Tests
		{
			name:           "Feature: Blog (Disabled)",
			url:            "/blog",
			handler:        blogHandler,
			envVars:        map[string]string{EnvEnableBlog: FeatureDisabled},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Feature: Blog (Enabled)",
			url:            "/blog",
			handler:        blogHandler,
			envVars:        map[string]string{EnvEnableBlog: FeatureEnabled},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Feature: Bookmarks (Disabled)",
			url:            "/bookmarks",
			handler:        bookmarksHandler,
			envVars:        map[string]string{EnvEnableBookmarks: FeatureDisabled},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Feature: Bookmarks (Enabled)",
			url:            "/bookmarks",
			handler:        bookmarksHandler,
			envVars:        map[string]string{EnvEnableBookmarks: FeatureEnabled},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			req := httptest.NewRequest(tt.method, tt.url, nil)
			if tt.method == "" {
				req.Method = http.MethodGet
			}
			rec := httptest.NewRecorder()

			tt.handler(rec, req)

			if tt.handler != nil {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rec.Body.String())
			}

			for k, v := range tt.expectedHeader {
				assert.Equal(t, v, rec.Header().Get(k))
			}
		})
	}
}
