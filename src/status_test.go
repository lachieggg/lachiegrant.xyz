package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetStatusContent verifies that getStatusContent returns non-empty HTML
// by executing htop and merging it with the neofetch output.
// This test requires htop and the neofetch.html file to be present.
func TestGetStatusContent(t *testing.T) {
	// Skip if required files don't exist (integration test dependency)
	if _, err := os.Stat(nfHtmlPath); os.IsNotExist(err) {
		t.Skipf("Skipping: required file %s not found", nfHtmlPath)
	}

	s, err := getStatusContent()
	assert.Nil(t, err)
	assert.NotEmpty(t, s, "status content should not be empty")
}
