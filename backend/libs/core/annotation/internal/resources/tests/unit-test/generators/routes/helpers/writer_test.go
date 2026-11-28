package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path/filepath"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/routes/helpers"
)

var _ = Describe("DiscoverWithTier", func() {
	It("omits route-limit entries for unlimited routes", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "protected/record.go", "package protected\nimport \"net/http\"\n// @Routes /records\ntype RecordController interface {\n// @GET /{gid}\nGet(http.ResponseWriter,*http.Request)\n}\n")
		groups, err := helpers.DiscoverWithTier(root, "/school", true, false)
		Expect(err).NotTo(HaveOccurred())
		out := filepath.Join(root, "out.txt")
		Expect(helpers.WriteRoutes("example.org/demo", "school", groups, out)).To(Succeed())
		content := fixture.Read(out)
		Expect(content).To(ContainSubstring("func RouteLimits() map[string]int"))
		Expect(content).NotTo(MatchRegexp(`"GET /school/protected/records/\{gid\}"\s*:`))
	})
})
