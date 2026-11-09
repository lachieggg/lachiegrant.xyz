package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFile_Success(t *testing.T) {
	// Create a temporary file with known content
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	expectedContent := "Hello!\nWorld!"

	err := os.WriteFile(tmpFile, []byte(expectedContent), 0644)
	require.NoError(t, err)

	content, err := ReadFile(tmpFile)

	require.NoError(t, err)
	assert.Equal(t, expectedContent, content)
}

func TestReadFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty")

	err := os.WriteFile(tmpFile, []byte{}, 0644)
	require.NoError(t, err)

	content, err := ReadFile(tmpFile)

	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestReadFile_FileNotFound(t *testing.T) {
	content, err := ReadFile("/bad/path/file.txt")

	assert.Error(t, err)
	assert.Empty(t, content)
	assert.True(t, os.IsNotExist(err))
}

func TestReadFile_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	// Attempting to read a directory should fail
	content, err := ReadFile(tmpDir)

	assert.Error(t, err)
	assert.Empty(t, content)
}

func TestReadFile_SpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "special-characters")
	expectedContent := "Unicode: 日本語 🎉"

	err := os.WriteFile(tmpFile, []byte(expectedContent), 0644)
	require.NoError(t, err)

	content, err := ReadFile(tmpFile)

	require.NoError(t, err)
	assert.Equal(t, expectedContent, content)
}
