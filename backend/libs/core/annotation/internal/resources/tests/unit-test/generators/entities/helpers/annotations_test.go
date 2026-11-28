package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/entities/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Generate", func() {
	It("rejects explicit accessors on persistent entities", func() {
		root := GinkgoT().TempDir()
		out := GinkgoT().TempDir()
		fixture.Write(root,"record.go","package entity\n// @Entity demo.records\n// @Getter\ntype Record struct{ Name string }\n")
		Expect(helpers.Generate(root,out)).To(MatchError(ContainSubstring("@Entity cannot request builders or accessors")))
	})
})
