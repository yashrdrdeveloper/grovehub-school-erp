package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path/filepath"
	"strings"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("ValidateImplementation", func() {
	It("rejects implementation structs without a database field", func() {
		root := searchFixture(strings.Replace(searchTestImpl, "db *gorm.DB", "other *gorm.DB", 1), searchTestConfig)
		specs, err := helpers.Discover(root, "grovehub", "demo")
		Expect(err).NotTo(HaveOccurred())
		Expect(helpers.ValidateImplementation(filepath.Join(root, "impl"), specs)).To(MatchError(ContainSubstring("needs a db field")))
	})
})
