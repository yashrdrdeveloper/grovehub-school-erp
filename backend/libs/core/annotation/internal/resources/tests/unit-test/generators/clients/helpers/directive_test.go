package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/clients/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover", func() {
	It("rejects duplicate proxy bindings", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "client.go", "package client\n// @Client clientImpl\ntype Client interface {\n// @Proxy service.Get\n// @Proxy service.Other\nGet(id string) string\n}\n")
		_, err := helpers.Discover(root, "")
		Expect(err).To(MatchError(ContainSubstring("duplicate @Proxy")))
	})
})
