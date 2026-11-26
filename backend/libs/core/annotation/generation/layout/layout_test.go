package layout_test

import (
	"os"
	"path/filepath"
	"testing"

	"grovehub/libs/core/annotation/generation/layout"
)

func TestReadModulePath(t *testing.T) {
	root := t.TempDir()
	if _, err := layout.ReadModulePath(root); err == nil {
		t.Fatal("missing module file accepted")
	}
	filename := filepath.Join(root, "go.mod")
	for _, entry := range []struct{ source, want string }{
		{"module grovehub\n", "grovehub"},
		{"// comment\nmodule example.org/internal-platform\n", "example.org/internal-platform"},
		{"go 1.27.0\n", ""},
	} {
		if err := os.WriteFile(filename, []byte(entry.source), 0600); err != nil {
			t.Fatal(err)
		}
		got, err := layout.ReadModulePath(root)
		if entry.want == "" {
			if err == nil {
				t.Fatal("missing module directive accepted")
			}
			continue
		}
		if err != nil || got != entry.want {
			t.Fatalf("got %q (%v), want %q", got, err, entry.want)
		}
	}
}

func TestResolveImport(t *testing.T) {
	root := t.TempDir()
	modulePath := "example.org/internal-platform"
	resolved, err := layout.ResolveImport(root, modulePath, modulePath+"/services/demo/dto")
	if err != nil || resolved != filepath.Join(root, "services", "demo", "dto") {
		t.Fatalf("unexpected import resolution: %q (%v)", resolved, err)
	}
	for _, invalid := range []string{"grovehub/services/demo/dto", modulePath + "/../escape", modulePath + "-other/services/demo"} {
		if _, err := layout.ResolveImport(root, modulePath, invalid); err == nil {
			t.Fatalf("invalid or foreign import accepted: %q", invalid)
		}
	}
}
