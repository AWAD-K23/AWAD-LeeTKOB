// Package localenv reads optional local settings without printing their values.
package localenv

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Load supports starting commands from the repository root or backend.
// An explicitly supplied environment variable always takes precedence.
func Load() error {
	for _, path := range []string{".env", filepath.Join("..", ".env")} {
		file, err := os.Open(path)

		if os.IsNotExist(err) {
			continue
		}

		if err != nil {
			return fmt.Errorf("open %s: %w", path, err)
		}

		defer file.Close()

		// Ignore blank lines and comments while reading local settings.
		scanner := bufio.NewScanner(file)
		for line := 1; scanner.Scan(); line++ {
			text := strings.TrimSpace(scanner.Text())

			if text == "" || strings.HasPrefix(text, "#") {
				continue
			}

			key, value, ok := strings.Cut(text, "=")
			key, value = strings.TrimSpace(key), strings.TrimSpace(value)

			if !ok || key == "" || strings.ContainsAny(key, " \t") {
				return fmt.Errorf(
					"%s:%d: invalid environment assignment", path, line,
				)
			}

			// Quoted values can contain spaces or assignment characters.
			if strings.HasPrefix(value, `"`) {
				value, err = strconv.Unquote(value)

				if err != nil {
					return fmt.Errorf("%s:%d: invalid quoted value", path, line)
				}
			} else if strings.HasPrefix(value, "'") {
				if len(value) < 2 || !strings.HasSuffix(value, "'") {
					return fmt.Errorf("%s:%d: invalid quoted value", path, line)
				}

				value = value[1 : len(value)-1]
			}

			// Preserve settings supplied by the shell or deployment environment.
			if _, exists := os.LookupEnv(key); !exists {
				if err := os.Setenv(key, value); err != nil {
					return fmt.Errorf("%s:%d: invalid environment name", path, line)
				}
			}
		}

		return scanner.Err()
	}

	// A local file is optional when settings are provided by the environment.
	return nil
}
