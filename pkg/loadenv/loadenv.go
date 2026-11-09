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
		if strings.HasPrefix(line, "#") {
			// Skip comments
			continue
		}
		if len(strings.TrimSpace(line)) == 0 {
			// Skip empty lines
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format in .env file: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = unquote(value)

		os.Setenv(key, value)
	}

	return nil
}

// unquote strips a single matching pair of surrounding single or double
// quotes from a .env value, as most .env conventions allow (e.g.
// SITE_NAME="My Site"). Values aren't required to be quoted -
// TrimSpace already preserves internal spaces - but a quoted value shouldn't
// end up with the literal quote characters in it. Left unfixed, those
// literal quote characters get percent-encoded (%22) wherever the value is
// used in an HTML attribute (see homeImage() in src/handlers.go), silently
// breaking the resulting URL.
func unquote(value string) string {
	if len(value) < 2 {
		return value
	}
	first, last := value[0], value[len(value)-1]
	if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
		return value[1 : len(value)-1]
	}
	return value
}
