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

func TestPickHomeImageName(t *testing.T) {
	tests := []struct {
		name     string
		pictures string
		unset    bool
		want     []string // acceptable outputs (a random pick among several)
	}{
		{
			name:     "single image",
			pictures: "beauty.png",
			want:     []string{"beauty.png"},
		},
		{
			name:     "picks one of several, trimmed",
			pictures: "a.png, b.png ,c.png",
			want:     []string{"a.png", "b.png", "c.png"},
		},
		{
			name:     "empty string falls back to placeholder",
			pictures: "",
			want:     []string{defaultHomeImage},
		},
		{
			name:  "unset env falls back to placeholder",
			unset: true,
			want:  []string{defaultHomeImage},
		},
		{
			name:     "blank entries between commas are ignored",
			pictures: "a.png,,  ,b.png",
			want:     []string{"a.png", "b.png"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.unset {
				os.Unsetenv(EnvPictures)
			} else {
				os.Setenv(EnvPictures, tt.pictures)
				defer os.Unsetenv(EnvPictures)
			}

			got := pickHomeImageName()
			assert.Contains(t, tt.want, got)
		})
	}
}

// TestHomeImage covers the layer on top of pickHomeImageName: preferring a
// pre-published, content-hash-versioned URL (see publishHomeImages in
// assets.go) and falling back to a plain /images/<name> path only when
// that image wasn't published - the exact fallback that would have masked
// the stale-Cloudflare-cache bug that hit the footer logo (rpi.png) if it
// had been in place, but must still degrade safely rather than 500 the
// homepage if a future image is somehow missing from the version map.
func TestHomeImage(t *testing.T) {
	t.Run("uses the pre-published versioned URL when available", func(t *testing.T) {
		os.Setenv(EnvPictures, "beauty.jpg")
		defer os.Unsetenv(EnvPictures)
		homeImageVersions["beauty.jpg"] = "/images/beauty.deadbeef1234.jpg"
		defer delete(homeImageVersions, "beauty.jpg")

		assert.Equal(t, "/images/beauty.deadbeef1234.jpg", homeImage())
	})

	t.Run("falls back to a plain path when not published", func(t *testing.T) {
		os.Setenv(EnvPictures, "unpublished.jpg")
		defer os.Unsetenv(EnvPictures)
		delete(homeImageVersions, "unpublished.jpg")

		assert.Equal(t, "/images/unpublished.jpg", homeImage())
	})
}

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

		// Code Page Handler Tests
		{
			name:    "Page: Code (renders both platform links)",
			url:     "/code",
			handler: codeHandler,
			envVars: map[string]string{
				EnvGithubURL: "https://github.com/lachieggg",
				EnvGitlabURL: "https://gitlab.com/lachie-gitlab",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Page: Code (renders even with URLs unset)",
			url:            "/code",
			handler:        codeHandler,
			envVars:        map[string]string{EnvGithubURL: "", EnvGitlabURL: ""},
			expectedStatus: http.StatusOK,
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

func TestCodeHandlerLinksToBothPlatforms(t *testing.T) {
	os.Setenv(EnvGithubURL, "https://github.com/lachieggg")
	os.Setenv(EnvGitlabURL, "https://gitlab.com/lachie-gitlab")
	defer os.Unsetenv(EnvGithubURL)
	defer os.Unsetenv(EnvGitlabURL)

	req := httptest.NewRequest(http.MethodGet, "/code", nil)
	rec := httptest.NewRecorder()

	codeHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `href="https://github.com/lachieggg"`)
	assert.Contains(t, body, `href="https://gitlab.com/lachie-gitlab"`)
	assert.Contains(t, body, "GitHub")
	assert.Contains(t, body, "GitLab")
}

func TestNotFoundIsThemed(t *testing.T) {
	rec := httptest.NewRecorder()
	indexHandler(rec, httptest.NewRequest(http.MethodGet, "/no-such-page", nil))

	assert.Equal(t, http.StatusNotFound, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "404")
	assert.Contains(t, body, "data-theme") // the shared header's theme script, i.e. full site chrome
	assert.Contains(t, body, "Return home")
}
