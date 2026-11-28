package testgen_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/generators/testing"
	"os"
	"path/filepath"
)

var _ = Describe("Unannotated test declarations", func() {
	It("ignores types and functions without testing annotations", func() {
		root := GinkgoT().TempDir()
		dir := filepath.Join(root, "tests", "unit-test", "plain")
		Expect(os.MkdirAll(dir, 0755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(dir, "plain_test.go"), []byte("package plain_test\ntype Plain struct{}\nfunc helper() {}"), 0644)).To(Succeed())
		Expect(testgen.GenerateAll(root, root)).To(Succeed())
	})
})
