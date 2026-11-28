package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("Discover", func() {
	It("requires a static coresearch.Config literal", func() {
		root := searchFixture(searchTestImpl, strings.Replace(searchTestConfig, "coresearch.Config{", "OtherConfig{", 1))
		specs, err := helpers.Discover(root, "grovehub", "demo")
		Expect(err).To(MatchError(ContainSubstring("must be coresearch.Config literal")))
		Expect(specs).To(BeNil())
	})
})
