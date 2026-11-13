package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidationCommentsRequireMatchingExecutableTags(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dto.go")
	broken := "package dto\ntype Input struct {\n // @NotBlank\n Name string\n}\n"
	if err := os.WriteFile(path, []byte(broken), 0600); err != nil { t.Fatal(err) }
	if err := CheckSource(path); err == nil || !strings.Contains(err.Error(), "notblank") {
		t.Fatalf("expected missing tag to fail, got %v", err)
	}
	valid := "package dto\ntype Input struct {\n // @NotBlank\n // @Size max=100\n Name string `validate:\"notblank,maxrunes=100\"`\n}\n"
	if err := os.WriteFile(path, []byte(valid), 0600); err != nil { t.Fatal(err) }
	if err := CheckSource(path); err != nil { t.Fatalf("expected matching validation tags: %v", err) }
}

func TestStrategyAnnotationAcceptedOnInterfaceAndImplementation(t *testing.T) {
    path := filepath.Join(t.TempDir(), "strategy.go")
    source := "package strategies\n// @Strategy\ntype CreateStrategy interface { Execute() }\n// @Strategy\ntype Creator struct{}\n"
    if err := os.WriteFile(path, []byte(source), 0600); err != nil { t.Fatal(err) }
    if err := CheckSource(path); err != nil { t.Fatalf("strategy interface and implementation should be accepted: %v", err) }
}
