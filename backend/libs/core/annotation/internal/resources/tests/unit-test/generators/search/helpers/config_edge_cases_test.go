package helpers_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Search configuration ownership and validation",func(){
	It("requires a literal sort field and matching cursor expression",func(){
		for _,test:=range []struct{old,new,reason string}{
			{"SortFields: map[string]coresearch.SortField{\"gid\": {Expression: \"p.gid\"}}","SortFields: other","expected static map literal"},
			{"SortFields: map[string]coresearch.SortField{\"gid\": {Expression: \"p.gid\"}}","SortFields: map[string]coresearch.SortField{\"gid\": other}","must use a literal"},
			{"SortFields: map[string]coresearch.SortField{\"gid\": {Expression: \"p.gid\"}}","SortFields: map[string]coresearch.SortField{\"gid\": {Expression: other}}","expected constant string"},
			{"SortFields: map[string]coresearch.SortField{\"gid\": {Expression: \"p.gid\"}}","SortFields: map[string]coresearch.SortField{\"gid\": {Expression: \"p.other\"}}","must match CursorExpression"},
		}{
			root:=searchFixture(strings.Replace(searchTestImpl,"alias=p","alias=p\n// @Cursor field=gid",1),strings.Replace(searchTestConfig,test.old,test.new,1))
			_,err:=helpers.Discover(root,"grovehub","demo")
			Expect(err).To(MatchError(ContainSubstring(test.reason)),test.new)
		}
	})
	It("resolves local configuration and resource-owned internal repository imports",func(){
		root:=searchFixture(strings.Replace(searchTestImpl,"config=search.PersonConfig","config=PersonConfig",1),searchTestConfig)
		fixture.Write(root,"impl/config.go",strings.Replace(searchTestConfig,"package search","package impl",1))
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred())
		out:=filepath.Join(root,"out.go")
		Expect(helpers.Write(out,"grovehub",specs)).To(Succeed())
		Expect(fixture.Read(out)).To(ContainSubstring("query, PersonConfig"))
		parent:=GinkgoT().TempDir()
		dir:=filepath.Join(parent,"internal","app","repository","person")
		fixture.Write(dir,"person.go",searchTestInterface)
		fixture.Write(dir,"impl/person.go",searchTestImpl)
		fixture.Write(dir,"helpers/search/person.go",searchTestConfig)
		fixture.Write(dir,"search/ignored.go","package other")
		fixture.Write(dir,"ignored_test.go","invalid source")
		fixture.Write(dir,"impl/ignored_test.go","invalid source")
		specs,err=helpers.Discover(dir,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred())
		out=filepath.Join(parent,".generated","internal","app","repository","person","impl","zz_search.gen.go")
		Expect(helpers.Write(out,"grovehub",specs)).To(Succeed())
		content:=fixture.Read(out)
		Expect(content).To(ContainSubstring("grovehub/services/demo/internal/app/repository/person/helpers/search"))
		Expect(content).To(ContainSubstring("grovehub/services/demo/internal/app/entity"))
	})
	It("rejects malformed implementation and config package source",func(){
		root:=searchFixture(searchTestImpl,searchTestConfig)
		fixture.Write(root,"impl/a.go","package impl\nfunc broken(")
		_,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(HaveOccurred())
		root=searchFixture(searchTestImpl,searchTestConfig)
		fixture.Write(root,"search/a.go","not a package")
		_,err=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(HaveOccurred())
		root=searchFixture(searchTestImpl,"package search\nvar broken =")
		_,err=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(HaveOccurred())
	})
	It("skips nonrepository declarations and unrelated interface methods",func(){
		root:=searchFixture(searchTestImpl,searchTestConfig)
		fixture.Write(root,"a.go","package repository\nconst x=1\ntype PersonRepository struct{}\ntype Plain interface{}\n")
		fixture.Write(root,"person.go",strings.Replace(searchTestInterface,"Search(ctx","Other\nOtherMethod()\nSearch(ctx",1))
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred());Expect(specs).To(HaveLen(1))
	})
	It("reports invalid glob patterns at the input boundary",func(){
		root:=filepath.Join(GinkgoT().TempDir(),"[")
		_,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(HaveOccurred())
		Expect(helpers.ValidateImplementation(root,nil)).NotTo(Succeed())
	})
})
