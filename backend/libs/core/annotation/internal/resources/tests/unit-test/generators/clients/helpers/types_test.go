package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/clients/helpers"
)

var _ = Describe("Discover", func() {
	It("rejects a request-line signature with a non-context first parameter", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "client.go", "package client\nimport \"context\"\n// @Client clientImpl\ntype Client interface {\n// @RequestLine GET /records\nGet(ctx string, input Request) (string, error)\n}\n")
		routes := fixture.Write(GinkgoT().TempDir(), "routes.go", "package routes\nfunc register(m M){m.Handle(\"GET /records\",nil)}")
		_, err := helpers.Discover(root, routes)
		Expect(err).To(MatchError(ContainSubstring("first @RequestLine parameter must be context.Context")))
	})
})
