package inventory_test

import (
    "os"
    "path/filepath"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing/inventory"
)

var _ = Describe("Inspect", func() {
    It("counts public and private mirrored tests without counting test support as production", func() {
        root := GinkgoT().TempDir()
        sources := map[string]string{
            "demo/demo.go": "package demo\nfunc Run() {}\n",
            "demo/internal/src/impl/run.go": "package impl\nfunc Execute() {}\n",
            "tests/unit-test/demo/demo_test.go": "package demo_test\n",
            "demo/internal/resources/tests/unit-test/impl/run_test.go": "package impl_test\n",
            "demo/internal/resources/tests/helpers/fixture.go": "package helpers\nfunc Fixture() {}\n",
        }

        for relative, source := range sources {
            path := filepath.Join(root, filepath.FromSlash(relative))
            Expect(os.MkdirAll(filepath.Dir(path), 0755)).To(Succeed())
            Expect(os.WriteFile(path, []byte(source), 0600)).To(Succeed())
        }

        report, err := inventory.Inspect(root)

        Expect(err).NotTo(HaveOccurred())
        Expect(report.SourceFiles).To(Equal(2))
        Expect(report.Functions).To(Equal(2))
        Expect(report.MirroredFunctions).To(Equal(2))
        Expect(report.TestFiles).To(Equal(2))
        Expect(report.Missing).To(BeEmpty())
    })
})
