package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/routes/helpers"
)

var _ = Describe("DiscoverWithTier", func() {
	It("rejects a mismatched controller package", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "protected/record.go", "package public\nimport \"net/http\"\n// @Routes /records\ntype RecordController interface {\n// @GET /{gid}\nGet(http.ResponseWriter,*http.Request)\n}\n")
		groups, err := helpers.DiscoverWithTier(root, "/school", true, false)
		Expect(err).To(MatchError(ContainSubstring("expected controller package")))
		Expect(groups).To(BeNil())
	})
})
