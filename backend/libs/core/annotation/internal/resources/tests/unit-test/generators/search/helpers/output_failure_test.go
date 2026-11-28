package helpers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Generated output persistence",func(){
	It("preserves deterministic output and reports invalid output destinations",func(){
		root:=GinkgoT().TempDir()
		root=searchFixture(searchTestImpl,searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred())
		save:=func(path string)error{return helpers.Write(path,"grovehub",specs)}
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
