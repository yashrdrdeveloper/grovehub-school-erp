package clients_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/generators/clients"
	"os"
	"path/filepath"
)

var _ = Describe("Source directory read errors", func() {
	It("reports unreadable source directories", func() {
		if os.Geteuid() == 0 {
			Skip("root bypasses permission checks")
		}
		root := GinkgoT().TempDir()
		Expect(os.WriteFile(filepath.Join(root, "go.mod"), []byte("module grovehub\n"), 0644)).To(Succeed())
		dir := filepath.Join(root, "services", "demo", "client")
		Expect(os.MkdirAll(dir, 0755)).To(Succeed())
		Expect(os.Chmod(dir, 0000)).To(Succeed())
		DeferCleanup(func() { _ = os.Chmod(dir, 0755) })
		Expect(clients.GenerateAll(root)).NotTo(Succeed())
	})
})
