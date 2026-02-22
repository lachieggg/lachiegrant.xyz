package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	// Test valid file creation
	t.Run("Valid Log File", func(t *testing.T) {
		initLogger(logPath)
		assert.NotNil(t, logger, "Logger should be initialized")

		logger.Println("Test message")

		// Verify file gets written
		data, err := os.ReadFile(logPath)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "Test message")
	})

	// Test invalid file creation (e.g. read only directory)
	t.Run("Invalid Log Path Gracefully Bounds to Stdout", func(t *testing.T) {
		invalidPath := filepath.Join(tempDir, "bad_dir", "test.log")
		initLogger(invalidPath)
		assert.NotNil(t, logger, "Logger should still dynamically initialize via Stdout")
	})

	// Test timing outputs
	t.Run("Timing Verification Context", func(t *testing.T) {
		initLogger(logPath)
		assert.Equal(t, log.LstdFlags, logger.Flags(), "Standard flags must always apply via fallback context")
	})
}
