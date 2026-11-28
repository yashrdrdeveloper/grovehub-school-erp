package testgen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateDirReportsInvalidDirectory(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "blocked")
	if err := os.WriteFile(path, nil, 0600); err != nil { t.Fatal(err) }
	if err := generateDir(root, path, "blocked"); err == nil { t.Fatal("expected directory read error") }
}
