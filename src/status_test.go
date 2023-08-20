package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStatusHandler
func TestStatusHandler(t *testing.T) {
	s, err := getStatusContent()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(s))
}

// WriteStringToFile writes the provided content to a file with the given name.
func WriteStringToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		return err
	}

	return file.Sync()
}
