package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindModuleDir(t *testing.T) {
	root := t.TempDir()
	module := filepath.Join(root, "project")
	child := filepath.Join(module, "services", "example", "internal", "tests", "mocks")
	if err := os.MkdirAll(child, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.test/project\n"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := findModuleDir(child)
	if err != nil || got != module {
		t.Fatalf("module = %q, error = %v", got, err)
	}
}
