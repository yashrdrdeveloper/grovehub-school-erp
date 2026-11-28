package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/clients/helpers"
)

var _ = Describe("Discover", func() {
	It("rejects noncanonical HTTP verbs before reading route metadata", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "client.go", "package client\n// @Client clientImpl\ntype Client interface {\n// @RequestLine TRACE /records\nGet(id string) string\n}\n")
		_, err := helpers.Discover(root, "")
		Expect(err).To(MatchError(ContainSubstring("unsupported @RequestLine verb")))
	})
})
