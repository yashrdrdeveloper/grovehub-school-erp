package clients_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/clients"
	sourceFixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Client resource output lifecycle",func(){
	It("writes to the implementation-owned overlay, removes stale outputs and reports blocked outputs",func(){
		root:=GinkgoT().TempDir()
		owner:=filepath.Join(root,"services","demo")
		source:="package record\n// @Client recordClient\ntype Client interface {\n// @Proxy service.Read\nRead(id string) string\n}"
		path:=sourceFixture.Write(owner,"client/record/client.go",source)
		Expect(os.MkdirAll(filepath.Join(owner,"internal","app","client","record","impl"),0755)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(owner,"client","impl"),0755)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(owner,"client",".hidden"),0755)).To(Succeed())
		Expect(clients.GenerateAll(root)).To(Succeed())
		out:=filepath.Join(owner,".generated","internal","app","client","record","impl","zz_clients.gen.go")
		Expect(sourceFixture.Read(out)).To(ContainSubstring("c.service.Read(id)"))
		Expect(os.Remove(out)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(out,"child"),0755)).To(Succeed())
		Expect(clients.GenerateAll(root)).NotTo(Succeed())
		sourceFixture.Write(owner,"client/record/client.go","package record\ntype Plain interface{}")
		Expect(clients.GenerateAll(root)).NotTo(Succeed())
		Expect(os.RemoveAll(out)).To(Succeed())
		Expect(clients.GenerateAll(root)).To(Succeed())
		Expect(path).To(BeAnExistingFile())
	})
	It("uses service-local internal route metadata for request-line clients",func(){
		root:=GinkgoT().TempDir()
		owner:=filepath.Join(root,"services","demo")
		sourceFixture.Write(owner,"client/client.go","package client\nimport \"context\"\n// @Client impl\ntype Client interface{\n// @RequestLine GET /protected/record\nRead(ctx context.Context, req Request)(Result,error)\n}")
		Expect(os.MkdirAll(filepath.Join(owner,"internal","app","routes"),0755)).To(Succeed())
		sourceFixture.Write(owner,".generated/internal/app/routes/zz_routes.gen.go","package routes\nfunc register(){mux.Handle(\"GET /demo/protected/record\",nil)}")
		Expect(clients.GenerateAll(root)).To(Succeed())
		Expect(sourceFixture.Read(filepath.Join(owner,".generated","client","impl","zz_clients.gen.go"))).To(ContainSubstring("/demo/protected/record"))
	})
})
