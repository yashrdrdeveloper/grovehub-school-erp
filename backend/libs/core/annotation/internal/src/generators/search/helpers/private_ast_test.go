package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindConfigReportsInvalidDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := findConfig(path, "Config"); err == nil {
		t.Fatal("expected directory read error")
	}
}

func TestFindConfigSkipsUninitializedVariables(t *testing.T) {
	dir := t.TempDir()
	source := "package search\nvar Other int\nvar Config = coresearch.Config{}\n"
	if err := os.WriteFile(filepath.Join(dir, "config.go"), []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	if config, err := findConfig(dir, "Config"); err != nil || config == nil {
		t.Fatalf("expected config literal, got %v, %v", config, err)
	}
}
