package search_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/generators/search"
	"os"
	"path/filepath"
)

var _ = Describe("GenerateAll", func() {
	It("reports a cyclic repository symlink", func() {
		root := GinkgoT().TempDir()
		Expect(os.WriteFile(filepath.Join(root, "go.mod"), []byte("module grovehub\n"), 0644)).To(Succeed())
		dir := filepath.Join(root, "services", "demo")
		Expect(os.MkdirAll(dir, 0755)).To(Succeed())
		Expect(os.Symlink("repository", filepath.Join(dir, "repository"))).To(Succeed())
		Expect(search.GenerateAll(root)).To(HaveOccurred())
	})
})
