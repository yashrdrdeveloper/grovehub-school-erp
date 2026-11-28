package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/clients/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover", func() {
	It("prefers the exact registered route over another suffix match", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "client.go", "package client\nimport \"context\"\n// @Client clientImpl\ntype Client interface {\n// @RequestLine GET /records\nGet(ctx context.Context, input Request) (string, error)\n}\n")
		routes := fixture.Write(GinkgoT().TempDir(),"routes.go","package routes\nfunc register(m M){m.Other(1); helper(1);m.Handle(\"GET /records\",nil);m.Handle(\"GET /other/records\",nil)}")
		proxies, err := helpers.Discover(root, routes)
		Expect(err).NotTo(HaveOccurred())
		output := filepath.Join(GinkgoT().TempDir(),"client.go")
		Expect(helpers.WriteProxies(proxies,output)).To(Succeed())
		Expect(fixture.Read(output)).To(ContainSubstring(`"GET", "/records"`))
		Expect(fixture.Read(output)).NotTo(ContainSubstring("/other/records"))
	})
})
