package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("Discover", func() {
	It("rejects nonexistent config package", func() {
		root := searchFixture(strings.Replace(searchTestImpl, "search.PersonConfig", "missing.PersonConfig", 1), searchTestConfig)
		specs, err := helpers.Discover(root, "grovehub", "demo")
		Expect(err).To(MatchError(ContainSubstring("config package")))
		Expect(specs).To(BeNil())
	})
})
