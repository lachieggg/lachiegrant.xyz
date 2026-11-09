package loadenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	// Create a temporary .env file for testing
	envContent := []byte("VAR_1=value-1\nVAR_2=value-2\nVAR_3=value=3\nVAR_4=\"quoted value\"\nVAR_5='single-quoted'\n# This is a comment\n")
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
	value1 := os.Getenv("VAR_1")
	assert.Equal(t, "value-1", value1)

	value2 := os.Getenv("VAR_2")
	assert.Equal(t, "value-2", value2)

	value3 := os.Getenv("VAR_3")
	assert.Equal(t, "value=3", value3)

	// Quoted values should have the literal quote characters stripped -
	// left in place, they end up percent-encoded wherever the value is
	// later used in an HTML attribute/URL.
	value4 := os.Getenv("VAR_4")
	assert.Equal(t, "quoted value", value4)

	value5 := os.Getenv("VAR_5")
	assert.Equal(t, "single-quoted", value5)

	// Test if a comment was ignored
	valueComment := os.Getenv("# This is a comment")
	assert.Empty(t, valueComment)
}

func TestUnquote(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"double-quoted", `"beauty.png"`, "beauty.png"},
		{"single-quoted", "'beauty.png'", "beauty.png"},
		{"unquoted", "beauty.png", "beauty.png"},
		{"empty string", "", ""},
		{"single char", `"`, `"`},
		{"mismatched quotes left alone", `"beauty.png'`, `"beauty.png'`},
		{"internal quotes preserved", `"say ""hi"""`, `say ""hi""`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, unquote(tt.input))
		})
	}
}
