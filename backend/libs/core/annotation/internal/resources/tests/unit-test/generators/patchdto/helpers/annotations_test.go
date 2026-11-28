package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/patchdto/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover",func(){
	It("rejects duplicate patch directives",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"input.go","package dto\nimport core \"grovehub/libs/core/annotation\"\n// @PatchRequestDto\n// @PatchRequestDto\ntype Input struct {\nName core.Patch[string] `json:\"name\"`\n}\n")
		specs,pkg,err:=helpers.Discover(root)
		Expect(err).To(MatchError(ContainSubstring("duplicate @PatchRequestDto")))
		Expect(specs).To(BeNil())
		Expect(pkg).To(BeEmpty())
	})
})
