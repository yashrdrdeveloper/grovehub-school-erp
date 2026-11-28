package inventory_test

import (
	"os"
	"path/filepath"
	"grovehub/libs/core/annotation/internal/src/generators/testing/inventory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inspect", func() {
	It("returns source traversal failures", func() {
		if os.Geteuid() == 0 { Skip("root bypasses directory permissions") }
		root := GinkgoT().TempDir()
		dir := filepath.Join(root, "private")
		Expect(os.Mkdir(dir, 0700)).To(Succeed())
		Expect(os.Chmod(dir, 0000)).To(Succeed())
		DeferCleanup(func() { _ = os.Chmod(dir, 0700) })
		_, err := inventory.Inspect(root)
		Expect(err).To(HaveOccurred())
	})
	It("sorts missing mirrors by source path", func() {
		root := GinkgoT().TempDir()
		Expect(os.WriteFile(filepath.Join(root, "a.go"), []byte("package demo\nfunc A() {}"), 0600)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(root, "b.go"), []byte("package demo\nfunc B() {}"), 0600)).To(Succeed())
		report, err := inventory.Inspect(root)
		Expect(err).NotTo(HaveOccurred())
		Expect(report.Missing).To(HaveLen(2))
		Expect(report.Missing[0].Source).To(Equal("a.go"))
		Expect(report.Missing[1].Source).To(Equal("b.go"))
	})
})
