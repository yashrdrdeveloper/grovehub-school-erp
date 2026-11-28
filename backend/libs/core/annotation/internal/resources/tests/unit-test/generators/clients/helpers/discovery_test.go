package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/clients/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover", func() {
	It("rejects unnamed method arguments", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "client.go", "package client\n// @Client clientImpl\ntype Client interface {\n// @Proxy service.Get\nGet(string) string\n}\n")
		_, err := helpers.Discover(root, "")
		Expect(err).To(MatchError(ContainSubstring("all arguments must be named")))
	})
})

var _ = Describe("Client signature imports", func() {
	It("copies signature qualifiers and omits facade-only imports", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "client.go", "package client\nimport (\n\"context\"\n\"strings\"\nmodel \"example/models\"\n\"example/unused\"\n\"example/module/v2\"\n\"example/values\"\n_ \"example/init\"\n. \"example/types\"\n)\nvar _ = strings.TrimSpace\n// @Client clientImpl\ntype Client interface {\n// @Proxy service.Get\nGet(ctx context.Context, id model.ID, v module.Value, x values.Item) (model.Result,error)\n}\n")
		proxies, err := helpers.Discover(root, "")
		Expect(err).NotTo(HaveOccurred())
		output := filepath.Join(root, "generated.go")
		Expect(helpers.WriteProxies(proxies, output)).To(Succeed())
		generated := fixture.Read(output)
		for _, expected := range []string{"\"context\"", "model \"example/models\"", "\"example/module/v2\"", "\"example/values\"", ". \"example/types\""} {
			Expect(generated).To(ContainSubstring(expected))
		}
		for _, unused := range []string{"\"strings\"", "\"example/unused\"", "\"example/init\""} {
			Expect(generated).NotTo(ContainSubstring(unused))
		}
	})
})
