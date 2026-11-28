package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/clients/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover", func() {
	It("retains two distinct methods in a single client proxy", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "client.go", "package client\n// @Client clientImpl\ntype Client interface {\n// @Proxy service.Get\nGet(id string) string\n// @Proxy service.Delete\nDelete(id string)\n}\n")
		proxies, err := helpers.Discover(root, "")
		Expect(err).NotTo(HaveOccurred())
		Expect(proxies).To(HaveLen(1))
		output := filepath.Join(GinkgoT().TempDir(),"client.go")
		Expect(helpers.WriteProxies(proxies,output)).To(Succeed())
		Expect(fixture.Read(output)).To(ContainSubstring("return c.service.Get(id)"))
		Expect(fixture.Read(output)).To(ContainSubstring("c.service.Delete(id)"))
		Expect(fixture.Read(output)).NotTo(ContainSubstring("return c.service.Delete"))
	})
})
