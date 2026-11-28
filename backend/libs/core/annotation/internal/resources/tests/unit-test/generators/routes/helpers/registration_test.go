package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path/filepath"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/routes/helpers"
)

var _ = Describe("DiscoverWithTier", func() {
	It("derives a resource-specific package alias and registration name", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "records/protected/record.go", "package protected\nimport \"net/http\"\n// @Routes /records\ntype RecordController interface {\n// @GET /{gid}\nGet(http.ResponseWriter,*http.Request)\n}\n")
		groups, err := helpers.DiscoverWithTier(root, "/school", true, false)
		Expect(err).NotTo(HaveOccurred())
		out := filepath.Join(root, "out.txt")
		Expect(helpers.WriteRoutes("example.org/demo", "school", groups, out)).To(Succeed())
		content := fixture.Read(out)
		Expect(content).To(ContainSubstring("func registerProtectedRecordRoutes"))
		Expect(content).To(ContainSubstring(`recordsprotected "example.org/demo/services/school/controller/records/protected"`))
	})
})
