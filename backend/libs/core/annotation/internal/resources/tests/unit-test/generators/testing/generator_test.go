package testgen_test

import (

	"os"
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing"
)

var _ = Describe("FindModuleDir", func() {
	It("Find Module Dir", func() {
		t := GinkgoT()
	root := t.TempDir()
	module := filepath.Join(root, "project")
	child := filepath.Join(module, "services", "example", "internal", "resources", "tests", "mocks")
	Expect(os.MkdirAll(child, 0755)).To(Succeed())
	Expect(os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.test/project\n"), 0600)).To(Succeed())
	got, err := testgen.FindModuleDir(child)
	Expect(err).NotTo(HaveOccurred())
	Expect(got).To(Equal(module))

	})
})

var _ = Describe("GenerateAll", func() {
    It("discovers both public and capability-local test suites", func() {
        root := GinkgoT().TempDir()
        paths := []string{
            filepath.Join("tests", "unit-test", "demo"),
            filepath.Join("demo", "internal", "resources", "tests", "unit-test", "impl"),
        }
        for _, relative := range paths {
            dir := filepath.Join(root, relative)
            Expect(os.MkdirAll(dir, 0755)).To(Succeed())
            source := "package demo_test\n// @Testing suite=Demo\nfunc specs() {}\n"
            Expect(os.WriteFile(filepath.Join(dir, "demo_test.go"), []byte(source), 0600)).To(Succeed())
        }

        Expect(testgen.GenerateAll(root, root)).To(Succeed())

        for _, relative := range paths {
            output := filepath.Join(root, ".generated", relative, "zz_suite_test.go")
            data, err := os.ReadFile(output)
            Expect(err).NotTo(HaveOccurred())
            Expect(string(data)).To(ContainSubstring("func TestDemo("))
            Expect(string(data)).To(ContainSubstring("specs()"))
        }
    })
})
