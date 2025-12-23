package exporter

import (
	"fmt"
	"strings"

	"golang.org/x/tools/imports"
)

// FormatterError represents an error that occurred during code formatting.
// It includes the original error and the unformatted code with line numbers.
type FormatterError struct {
	OrigErr error
	Code    string
}

// Error returns a formatted error message with the unformatted code.
func (e *FormatterError) Error() string {
	sb := &strings.Builder{}
	sb.WriteString("Error formatting generated code: " + e.OrigErr.Error() + "\n")
	for n, line := range strings.Split(e.Code, "\n") {
		sb.WriteString(fmt.Sprintf("  %4d: %s\n", n+1, line))
	}
	return sb.String()
}

// Unwrap returns the original error.
func (e *FormatterError) Unwrap() error {
	return e.OrigErr
}

// formatCode formats the given Go code using the imports package.
// If formatting fails, it returns a FormatterError containing the original error and the unformatted code.
func formatCode(code string) (string, error) {
	formatted, err := imports.Process("code.go", []byte(code), nil)
	if err != nil {
		return "", &FormatterError{OrigErr: err, Code: code}
	}
	return string(formatted), nil
}
