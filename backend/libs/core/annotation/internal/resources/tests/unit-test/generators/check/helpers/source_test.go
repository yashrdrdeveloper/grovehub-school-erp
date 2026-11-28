package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"grovehub/libs/core/annotation/internal/src/generators/check/helpers"
)

var _ = Describe("IsMaintainedSource",func(){
	DescribeTable("selects Go production sources",func(path string,expected bool){Expect(helpers.IsMaintainedSource(path)).To(Equal(expected))},Entry("source","file.go",true),Entry("test","file_test.go",false),Entry("markdown","file.md",false))
})
var _ = Describe("CheckDirectory",func(){
	It("skips generated and vendor directories",func(){
		root:=GinkgoT().TempDir()
		for _,name:=range []string{"vendor",".generated","source"}{Expect(os.Mkdir(filepath.Join(root,name),0700)).To(Succeed())}
		entries,err:=os.ReadDir(root);Expect(err).NotTo(HaveOccurred())
		for _,entry:=range entries{if entry.Name()=="source"{Expect(helpers.CheckDirectory(entry)).To(Succeed());continue};Expect(helpers.CheckDirectory(entry)).To(MatchError(filepath.SkipDir))}
	})
})
