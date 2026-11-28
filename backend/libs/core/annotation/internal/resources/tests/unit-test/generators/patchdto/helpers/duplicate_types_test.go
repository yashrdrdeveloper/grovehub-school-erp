package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/patchdto/helpers"
)

var _ = Describe("PATCH source ownership", func() {
	It("rejects duplicate patch types across files", func() {
		root := GinkgoT().TempDir()
		for _, file := range []string{"a.go", "b.go"} {
			fixture.Write(root, file, "package dto\n// @PatchRequestDto\ntype Update struct{}")
		}
		_, _, err := helpers.Discover(root)
		Expect(err).To(MatchError(ContainSubstring("duplicate @PatchRequestDto")))
	})
})
