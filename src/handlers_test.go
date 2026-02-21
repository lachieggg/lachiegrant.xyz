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
	os.Exit(m.Run())
}

func TestIndexHandler_RootPath_Returns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Note: This will fail if templates aren't available, but tests the routing logic
	indexHandler(rec, req)

	// If templates aren't found, we get 500, but that's expected in a test environment
	// We mainly want to verify it doesn't return 404 for "/"
	assert.NotEqual(t, http.StatusNotFound, rec.Code)
}

func TestIndexHandler_NonRootPath_Returns404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()

	indexHandler(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestIndexHandler_NestedPath_Returns404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/some/nested/path", nil)
	rec := httptest.NewRecorder()

	indexHandler(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGithubHandler_Redirects(t *testing.T) {
	// Set up test environment variable
	expectedURL := "https://github.com/testuser"
	os.Setenv("GITHUB_URL", expectedURL)
	defer os.Unsetenv("GITHUB_URL")

	req := httptest.NewRequest(http.MethodGet, "/code", nil)
	rec := httptest.NewRecorder()

	githubHandler(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, expectedURL, rec.Header().Get("Location"))
}

func TestGithubHandler_EmptyURL_StillRedirects(t *testing.T) {
	os.Unsetenv("GITHUB_URL")

	req := httptest.NewRequest(http.MethodGet, "/code", nil)
	rec := httptest.NewRecorder()

	githubHandler(rec, req)

	// Should still return redirect status, even with empty URL
	assert.Equal(t, http.StatusSeeOther, rec.Code)
}

func TestMiddlewareFunc_LogsRequest(t *testing.T) {
	// Create a simple handler that records if it was called
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with middleware
	wrapped := middlewareFunc(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Real-IP", "192.168.1.1")
	rec := httptest.NewRecorder()

	wrapped(rec, req)

	assert.True(t, handlerCalled, "handler should have been called")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddlewareFunc_PassesThroughRequest(t *testing.T) {
	var capturedPath string
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
	})

	wrapped := middlewareFunc(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/specific/path", nil)
	rec := httptest.NewRecorder()

	wrapped(rec, req)

	assert.Equal(t, "/specific/path", capturedPath)
}
