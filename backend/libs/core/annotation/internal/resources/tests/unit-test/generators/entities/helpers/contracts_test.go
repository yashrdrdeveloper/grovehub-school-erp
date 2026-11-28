package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/entities/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Generate", func() {
	It("preserves each field name and type in generated builders", func() {
		root := GinkgoT().TempDir()
		out := GinkgoT().TempDir()
		fixture.Write(root,"record.go","package entity\n// @Builder\ntype Record struct{ First, Last string; Count int }\n")
		Expect(helpers.Generate(root,out)).To(Succeed())
		content:=fixture.Read(filepath.Join(out,"record_gombok.go"))
		Expect(content).To(ContainSubstring("WithFirst(value string)"))
		Expect(content).To(ContainSubstring("WithLast(value string)"))
		Expect(content).To(ContainSubstring("WithCount(value int)"))
	})
})
