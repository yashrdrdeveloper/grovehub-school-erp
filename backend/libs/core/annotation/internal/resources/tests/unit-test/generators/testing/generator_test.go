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
	child := filepath.Join(module, "services", "example", "internal", "tests", "mocks")
	Expect(os.MkdirAll(child, 0755)).To(Succeed())
	Expect(os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.test/project\n"), 0600)).To(Succeed())
	got, err := testgen.FindModuleDir(child)
	Expect(err).NotTo(HaveOccurred())
	Expect(got).To(Equal(module))

	})
})
