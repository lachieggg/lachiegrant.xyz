package loadenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	// Create a temporary .env file for testing
	envContent := []byte("TEST_VAR_1=value1\nTEST_VAR_2=value2\n# This is a comment\n")
	tmpfile, err := os.CreateTemp("", "testenv")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(envContent)
	assert.NoError(t, err)

	err = tmpfile.Close()
	assert.NoError(t, err)

	// Load environment variables from the temporary .env file
	err = LoadEnv(tmpfile.Name())
	assert.NoError(t, err)

	// Test if the environment variables were set correctly
	value1 := os.Getenv("TEST_VAR_1")
	assert.Equal(t, "value1", value1)

	value2 := os.Getenv("TEST_VAR_2")
	assert.Equal(t, "value2", value2)

	// Test if a comment was ignored
	valueComment := os.Getenv("# This is a comment")
	assert.Empty(t, valueComment)
}
