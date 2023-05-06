package loadenv

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnv loads an env file into the
// environment
func LoadEnv(envPath string) error {
	envFile, err := os.Open(filepath.Clean(envPath))
	if err != nil {
		return err
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			// Skip comments and empty lines
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Invalid format in .env file: %s", line)
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	return nil
}
