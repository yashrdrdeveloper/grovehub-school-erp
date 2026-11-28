package routes_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/routes"
	sourceFixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Internal service route output lifecycle",func(){
	It("generates internal resource routes, checks service prefixes and handles stale output",func(){
		root:=GinkgoT().TempDir()
		sourceFixture.Write(root,"go.mod","module example.com/project\n")
		owner:=filepath.Join(root,"services","demo")
		sourceFixture.Write(owner,"internal/app/routes/http.go","package routes\n// @Routes /demo\nfunc RegisterHTTP(){}")
		sourceFixture.Write(owner,"internal/app/controller/record/protected/record.go","package protected\n// @Routes /records\ntype RecordController interface{\n// @GET /\nRead(http.ResponseWriter,*http.Request)\n}")
		Expect(routes.GenerateAll(root)).To(Succeed())
		out:=filepath.Join(owner,".generated","internal","app","routes","zz_routes.gen.go")
		Expect(sourceFixture.Read(out)).To(ContainSubstring("/internal/app/controller/record/protected"))
		Expect(routes.GenerateAll(root)).To(Succeed())
		Expect(os.Remove(out)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(out,"child"),0755)).To(Succeed())
		Expect(routes.GenerateAll(root)).NotTo(Succeed())
		sourceFixture.Write(owner,"internal/app/controller/record/protected/record.go","package protected\ntype Plain interface{}")
		Expect(routes.GenerateAll(root)).NotTo(Succeed())
		Expect(os.RemoveAll(out)).To(Succeed())
		Expect(routes.GenerateAll(root)).To(Succeed())
		Expect(out).NotTo(BeAnExistingFile())
		sourceFixture.Write(owner,"internal/app/routes/http.go","package routes\n// @Routes /{bad}\nfunc RegisterHTTP(){}")
		Expect(routes.GenerateAll(root)).To(MatchError(ContainSubstring("cannot contain path variables")))
	})
	It("rejects a controller file and a looping restricted marker",func(){
		root:=GinkgoT().TempDir()
		sourceFixture.Write(root,"go.mod","module example.com/project\n")
		sourceFixture.Write(root,"services/demo/controller","file")
		Expect(routes.GenerateAll(root)).To(MatchError(ContainSubstring("controller must be a directory")))
		Expect(os.Symlink(".restricted-routes",filepath.Join(root,".restricted-routes"))).To(Succeed())
		Expect(routes.GenerateAll(root)).NotTo(Succeed())
	})
})
