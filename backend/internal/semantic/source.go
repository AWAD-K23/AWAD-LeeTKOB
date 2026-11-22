package semantic

import (
	"errors"
	"strings"
)

// stripComments blanks comments without moving source lines or damaging quoted
// literals. Ambiguous source fails closed rather than leaking comment text.
func stripComments(source string) (string, error) {
	// Translation-phase line splices can extend // comments or split delimiters.
	// Reject them until a line-mapped preprocessor is available.
	if strings.Contains(source, "\\\n") || strings.Contains(source, "\\\r\n") ||
		strings.Contains(source, "??/") {
		return "", errors.New("unsupported source continuation")
	}

	clean := []byte(source)
	blank := func(start, end int) {
		for j := start; j < end; j++ {
			if clean[j] != '\n' && clean[j] != '\r' {
				clean[j] = ' '
			}
		}
	}

	for i := 0; i < len(source); {
		switch {
		case strings.HasPrefix(source[i:], "//"):
			end := i + 2
			for end < len(source) && source[end] != '\n' {
				end++
			}
			blank(i, end)
			i = end

		case strings.HasPrefix(source[i:], "/*"):
			end := strings.Index(source[i+2:], "*/")
			if end < 0 {
				return "", errors.New("unclosed comment")
			}
			end += i + 4
			blank(i, end)
			i = end

		case strings.HasPrefix(source[i:], `R"`):
			open := strings.IndexByte(source[i+2:], '(')
			if open < 0 || open > 16 {
				return "", errors.New("invalid raw literal")
			}
			delimiter := source[i+2 : i+2+open]
			closing := ")" + delimiter + `"`
			body := i + 3 + open
			end := strings.Index(source[body:], closing)
			if end < 0 {
				return "", errors.New("unclosed raw literal")
			}
			i = body + end + len(closing)

		case source[i] >= '0' && source[i] <= '9':
			// Consume preprocessing numbers together. Apostrophes in hexadecimal
			// and decimal constants are separators, not character-literal quotes.
			i++
			for i < len(source) {
				ch := source[i]
				if numberPart(ch) || ch == '.' ||
					(ch == '\'' && i+1 < len(source) && numberPart(source[i+1])) {
					i++
				} else if (ch == '+' || ch == '-') && strings.ContainsRune("eEpP", rune(source[i-1])) {
					i++
				} else {
					break
				}
			}

		case source[i] == '"' || source[i] == '\'':
			quote := source[i]
			i++
			closed := false
			for i < len(source) {
				if source[i] == '\n' || source[i] == '\r' {
					return "", errors.New("newline in ordinary literal")
				}
				if source[i] == '\\' {
					i += 2
				} else if source[i] == quote {
					i++
					closed = true
					break
				} else {
					i++
				}
			}
			if !closed {
				return "", errors.New("unclosed literal")
			}

		default:
			i++
		}
	}

	return string(clean), nil
}

func numberPart(ch byte) bool {
	return ch >= '0' && ch <= '9' || ch >= 'a' && ch <= 'z' ||
		ch >= 'A' && ch <= 'Z' || ch == '_'
}
