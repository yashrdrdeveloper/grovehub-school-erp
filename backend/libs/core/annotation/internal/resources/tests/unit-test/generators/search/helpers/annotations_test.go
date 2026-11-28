package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("Discover", func() {
	It("rejects duplicate searchable options", func() {
		root := searchFixture(strings.Replace(searchTestImpl, "alias=p", "alias=p alias=q", 1), searchTestConfig)
		specs, err := helpers.Discover(root, "grovehub", "demo")
		Expect(err).To(MatchError(ContainSubstring("invalid or repeated @Searchable option")))
		Expect(specs).To(BeNil())
	})
})
