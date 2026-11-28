package search_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/search"
	sourceFixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Resource-owned search generation",func(){
	It("uses internal repository paths, validates implementations, and preserves maintained sources",func(){
		root,_:=searchFixture(GinkgoT(),searchTestInterface,searchTestImpl)
		owner:=filepath.Join(root,"services","demo")
		old:=filepath.Join(owner,"repository")
		dir:=filepath.Join(owner,"internal","app","repository","person")
		Expect(os.MkdirAll(filepath.Dir(dir),0755)).To(Succeed())
		Expect(os.Rename(old,dir)).To(Succeed())
		for _,name:=range []string{"constants","helpers","impl",".hidden"}{
			Expect(os.MkdirAll(filepath.Join(filepath.Dir(dir),name),0755)).To(Succeed())
		}
		Expect(search.GenerateAll(root)).To(Succeed())
		out:=filepath.Join(owner,".generated","internal","app","repository","person","impl","zz_search.gen.go")
		Expect(sourceFixture.Read(out)).To(ContainSubstring("/internal/app/entity"))
		Expect(os.Remove(out)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(out,"child"),0755)).To(Succeed())
		Expect(search.GenerateAll(root)).NotTo(Succeed())
		Expect(os.RemoveAll(out)).To(Succeed())
		sourceFixture.Write(dir,"impl/person.go","package impl\ntype PersonRepository struct{}")
		Expect(search.GenerateAll(root)).To(Succeed())
		Expect(out).NotTo(BeAnExistingFile())
		Expect(os.MkdirAll(filepath.Join(out,"child"),0755)).To(Succeed())
		Expect(search.GenerateAll(root)).NotTo(Succeed())
	})
	It("reports a repository path that is not a directory",func(){
		root:=GinkgoT().TempDir()
		sourceFixture.Write(root,"go.mod","module example.com/project\n")
		sourceFixture.Write(root,"services/demo/repository","file")
		Expect(search.GenerateAll(root)).NotTo(Succeed())
	})
})
