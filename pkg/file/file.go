package file

import "os"

// ReadFile reads the content of a file and returns it as a string.
func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
