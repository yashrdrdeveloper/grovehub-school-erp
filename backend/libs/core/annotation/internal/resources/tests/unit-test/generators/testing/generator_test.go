package testgen_test

import (
	"os"
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing"
)

var _ = Describe("FindModuleDir", func() {
	It("Find Module Dir", func() {
		t := GinkgoT()
	root := t.TempDir()
	module := filepath.Join(root, "project")
	child := filepath.Join(module, "services", "example", "internal", "tests", "mocks")
	if err := os.MkdirAll(child, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.test/project\n"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := testgen.FindModuleDir(child)
	if err != nil || got != module {
		t.Fatalf("module = %q, error = %v", got, err)
	}

	})
})
