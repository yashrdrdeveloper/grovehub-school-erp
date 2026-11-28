package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/clients/helpers"
)

var _ = Describe("Client directive boundaries", func() {
	It("rejects empty or multiword bindings", func() {
		for _, annotation := range []string{"@Client", "@Client one two"} {
			root := GinkgoT().TempDir()
			fixture.Write(root, "client.go", "package client\n// "+annotation+"\ntype Client interface{}")
			_, err := helpers.Discover(root, "")
			Expect(err).To(MatchError(ContainSubstring("expects one binding")))
		}
		root := GinkgoT().TempDir()
		fixture.Write(root, "client.go", "package client\n// @Client impl\ntype Client interface{\n// @RequestLine GET /records\n// @RequestLineExtra ignore\nRead(ctx context.Context, req Request)(Result,error)\n}")
		route := fixture.Write(root, "routes/routes.go", "package routes\nfunc register(){mux.Handle(\"GET /records\",nil)}")
		_, err := helpers.Discover(root, route)
		Expect(err).NotTo(HaveOccurred())
	})
})
