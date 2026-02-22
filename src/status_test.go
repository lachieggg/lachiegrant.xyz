package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCmd(t *testing.T) {
	// 1. Using a generic 'echo' command since aha or htop might not exist on the runner physically
	// It will attempt to run `echo q | echo 'localhost test' | aha`, failing gracefully if aha is missing.
	output := ExecuteCmd("echo 'localhost test'")

	// Either aha exists, or the string execution gracefully caught the internal echo regardless.
	assert.Contains(t, output, "localhost", "Should securely execute a generic host test output")
}

func TestStatusHandlerMissingFile(t *testing.T) {
	// Temporarily override nfHtmlPath
	origPath := nfHtmlPath
	nfHtmlPath = "./does-not-exist.html"
	defer func() { nfHtmlPath = origPath }()

	req := httptest.NewRequest(http.MethodGet, "http://localhost/status", nil)
	rec := httptest.NewRecorder()

	statusHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected Status code to gracefully fail due to missing neofetch file")
}

func TestStatusHandlerValidExecution(t *testing.T) {
	tempDir := t.TempDir()
	testNeofetch := tempDir + "/neofetch.html"
	dummyContent := "<html><body><h1>ok</h1></body></html>"
	err := os.WriteFile(testNeofetch, []byte(dummyContent), 0644)
	assert.NoError(t, err)

	// Temporarily override nfHtmlPath
	origPath := nfHtmlPath
	nfHtmlPath = testNeofetch
	defer func() { nfHtmlPath = origPath }()

	req := httptest.NewRequest(http.MethodGet, "http://localhost/status", nil)
	rec := httptest.NewRecorder()

	statusHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected Success because neofetch structure existed")
	assert.Contains(t, rec.Body.String(), "ok", "Expected dynamically mocked content inside response")
}
