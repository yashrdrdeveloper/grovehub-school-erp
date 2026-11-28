package parser_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/parser"
	"os"
	"path/filepath"
)

var _ = Describe("ScanFile", func() {
	It("reports Go parse and missing-file errors without returning partial occurrences", func() {
		path := filepath.Join(GinkgoT().TempDir(), "source.go")
		occurrences, err := parser.ScanFile(path)
		Expect(err).To(HaveOccurred())
		Expect(occurrences).To(BeNil())
		Expect(os.WriteFile(path, []byte("package demo\nfunc broken("), 0600)).To(Succeed())
		occurrences, err = parser.ScanFile(path)
		Expect(err).To(HaveOccurred())
		Expect(occurrences).To(BeNil())
	})
})
