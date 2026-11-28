package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/entities/helpers"
)

var _ = Describe("Entity output directory", func() {
	It("reports an existing file in place of the output directory", func() {
		Expect(helpers.Generate(GinkgoT().TempDir(), fixture.Write(GinkgoT().TempDir(), "blocked", "file"))).NotTo(Succeed())
	})
})
