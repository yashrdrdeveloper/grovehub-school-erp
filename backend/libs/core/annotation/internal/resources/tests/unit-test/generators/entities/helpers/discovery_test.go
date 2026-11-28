package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/entities/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Generate", func() {
	It("rejects private fields in generated builders", func() {
		root := GinkgoT().TempDir()
		out := GinkgoT().TempDir()
		fixture.Write(root,"record.go","package entity\n// @Builder\ntype Record struct{ private string }\n")
		Expect(helpers.Generate(root,out)).To(MatchError(ContainSubstring("unexported field private")))
	})
})
