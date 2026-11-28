package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
	"strings"
)

var _ = Describe("Discover", func() {
	It("rejects an invalid qualified config reference", func() {
		root := searchFixture(strings.Replace(searchTestImpl, "search.PersonConfig", "search.extra.PersonConfig", 1), searchTestConfig)
		_, err := helpers.Discover(root, "grovehub", "demo")
		Expect(err).To(MatchError(ContainSubstring("invalid @Searchable config")))
	})
	It("ignores variable declarations without initializers", func() {
		root := searchFixture(searchTestImpl, searchTestConfig+"\nvar Other int\n")
		_, err := helpers.Discover(root, "grovehub", "demo")
		Expect(err).NotTo(HaveOccurred())
	})
	It("rejects a config without named fields", func() {
		root := searchFixture(searchTestImpl, "package search\nvar PersonConfig = coresearch.Config{1}")
		_, err := helpers.Discover(root, "grovehub", "demo")
		Expect(err).To(HaveOccurred())
	})
})
