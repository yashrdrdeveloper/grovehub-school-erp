// Package output persists formatted generated Go sources consistently.
package output

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
)

// WriteGo validates generated syntax before modifying a destination. The
// description preserves generator-specific diagnostic context when provided.
func WriteGo(path string, source []byte, description string) error {
	formatted, err := format.Source(source)
	if err != nil {
		if description != "" {
			return fmt.Errorf("%s: %w", description, err)
		}
		return err
	}
	return WriteIfChanged(path, formatted)
}

// WriteIfChanged leaves unchanged files untouched and propagates read, parent
// directory and write errors. Existing generated Go source modes are retained.
func WriteIfChanged(path string, contents []byte) error {
	old, err := os.ReadFile(path)
	if err == nil && bytes.Equal(old, contents) {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, contents, 0644)
}
