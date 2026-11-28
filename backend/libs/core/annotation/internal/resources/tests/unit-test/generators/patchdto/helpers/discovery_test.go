package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/patchdto/helpers"
)

var _ = Describe("Discover", func() {
	It("rejects duplicate JSON names", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "input.go", "package dto\nimport core \"grovehub/libs/core/annotation\"\n// @PatchRequestDto\ntype Input struct {\nA core.Patch[string] `json:\"same\"`\nB core.Patch[int] `json:\"same\"`\n}\n")
		specs, pkg, err := helpers.Discover(root)
		Expect(err).To(MatchError(ContainSubstring("duplicate json field")))
		Expect(specs).To(BeNil())
		Expect(pkg).To(BeEmpty())
	})
})
