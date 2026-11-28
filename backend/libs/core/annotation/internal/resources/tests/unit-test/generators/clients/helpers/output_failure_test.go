package helpers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/clients/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Generated output persistence",func(){
	It("preserves deterministic output and reports invalid output destinations",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"client.go","package client\nimport ctx \"context\"\n// @Client clientImpl\ntype Client interface {\n// @Proxy svc.Read\nRead(c ctx.Context) string\n// @Proxy svc.Clear\nClear()\n}")
		proxies,err:=helpers.Discover(root,"")
		Expect(err).NotTo(HaveOccurred())
		save:=func(path string)error{return helpers.WriteProxies(proxies,path)}
		out:=GinkgoT().TempDir()
		path:=filepath.Join(out,"generated.go")
		Expect(save(path)).To(Succeed())
		first:=fixture.Read(path)
		Expect(first).NotTo(BeEmpty())
		Expect(save(path)).To(Succeed())
		Expect(fixture.Read(path)).To(Equal(first))
		blocked:=fixture.Write(out,"blocked","not a directory")
		Expect(save(filepath.Join(blocked,"generated.go"))).NotTo(Succeed())
		Expect(fixture.Read(blocked)).To(Equal("not a directory"))
		directory:=filepath.Join(out,"directory.go")
		Expect(os.Mkdir(directory,0755)).To(Succeed())
		Expect(save(directory)).NotTo(Succeed())
	})
})
