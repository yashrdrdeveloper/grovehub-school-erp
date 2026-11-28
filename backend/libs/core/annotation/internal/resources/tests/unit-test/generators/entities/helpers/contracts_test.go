package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path/filepath"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/entities/helpers"
)

var _ = Describe("Generate", func() {
	It("preserves each field name and type in generated builders", func() {
		root := GinkgoT().TempDir()
		out := GinkgoT().TempDir()
		fixture.Write(root, "record.go", "package entity\n// @Builder\ntype Record struct{ First, Last string; Count int }\n")
		Expect(helpers.Generate(root, out)).To(Succeed())
		content := fixture.Read(filepath.Join(out, "record_gombok.go"))
		Expect(content).To(ContainSubstring("WithFirst(value string)"))
		Expect(content).To(ContainSubstring("WithLast(value string)"))
		Expect(content).To(ContainSubstring("WithCount(value int)"))
	})
})
