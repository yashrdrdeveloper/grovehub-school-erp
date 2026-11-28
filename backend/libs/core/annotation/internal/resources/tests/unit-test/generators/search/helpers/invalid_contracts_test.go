package helpers_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Search generator contract validation",func(){
	It("rejects malformed or incompatible concrete methods",func(){
		for _,test:=range []struct{old,new,reason string}{
			{"(repository *PersonRepository) Search","Search","requires a concrete Search method"},
			{"*PersonRepository) Search","*PersonRepository[T]) Search","named repository receiver"},
			{"ctx context.Context, query coresearch.Query[string]","ctx context.Context","requires (context.Context"},
			{"coresearch.Query[string]","coresearch.Query[int]","unsupported signature"},
			{"return repository.searchGenerated(ctx, query)","return","must return"},
			{"return repository.searchGenerated(ctx, query)","return query","must call"},
			{"repository.searchGenerated(ctx, query)","repository.other(ctx, query)","must use searchGenerated"},
			{"repository.searchGenerated(ctx, query)","other.searchGenerated(ctx, query)","repository receiver"},
			{"repository.searchGenerated(ctx, query)","repository.searchGenerated(query, ctx)","forward ctx and query"},
			{"config=search.PersonConfig alias=p","config=search.PersonConfig","requires config="},
			{"config=search.PersonConfig alias=p","config=search.PersonConfig alias=p alias=q","invalid or repeated"},
			{"config=search.PersonConfig alias=p","config=search.PersonConfig alias=p entity=Other","entity must match"},
			{"config=search.PersonConfig alias=p","config=search.PersonConfig alias=p impl=Other","impl must match"},
			{"config=search.PersonConfig alias=p","config=bad-config alias=p","invalid @Searchable config"},
			{"config=search.PersonConfig alias=p","config=missing.PersonConfig alias=p","config package"},
			{"alias=p","alias=p\n// @Filterable wrong=gid","@Filterable"},
			{"alias=p","alias=p\n// @Sortable wrong=gid","@Sortable"},
			{"alias=p","alias=p\n// @Cursor wrong=gid","@Cursor"},
			{"alias=p","alias=p\n// @Sortable fields=gid\n// @Cursor field=other","field must be sortable"},
		}{
			root:=searchFixture(strings.Replace(searchTestImpl,test.old,test.new,1),searchTestConfig)
			_,err:=helpers.Discover(root,"grovehub","demo")
			Expect(err).To(MatchError(ContainSubstring(test.reason)),test.new)
		}
	})
	It("rejects missing, malformed and misplaced repository declarations",func(){
		for _,test:=range []struct{source,reason string}{
			{"package repository\nfunc broken(","expected"},
			{strings.Replace(searchTestInterface,"// @Repository entity=Person","// @Repository wrong=Person",1),"matching @Repository"},
			{"package repository\n// @Repository entity=Person\ntype PersonRepository interface { Other() }","must declare Search"},
			{strings.Replace(searchTestInterface,"Search(ctx","// @Searchable config=x alias=p\nSearch(ctx",1),"move @Searchable"},
			{"package repository\ntype Other interface{}","missing @Repository"},
			{strings.Replace(searchTestInterface,"coresearch.Query[string]","coresearch.Query[int]",1),"unsupported signature"},
		}{
			root:=searchFixture(searchTestImpl,searchTestConfig)
			fixture.Write(root,"person.go",test.source)
			_,err:=helpers.Discover(root,"grovehub","demo")
			Expect(err).To(MatchError(ContainSubstring(test.reason)))
		}
		root:=searchFixture(searchTestImpl,searchTestConfig)
		fixture.Write(root,"impl/duplicate.go",searchTestImpl)
		_,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(MatchError(ContainSubstring("duplicate @Searchable")))
	})
	It("rejects invalid configuration literals and inconsistent cursor expressions",func(){
		for _,test:=range []struct{old,new,reason string}{
			{"var PersonConfig = coresearch.Config","var OtherConfig = coresearch.Config","not found"},
			{"coresearch.Config{","OtherConfig{","must be coresearch.Config literal"},
			{"RelationFields: map[string]coresearch.RelationField{}","RelationFields: missing","expected static map literal"},
			{"\"gid\": {Expression: \"p.gid\"}","1: {Expression: \"p.gid\"}","expected constant string"},
			{"CursorExpression: \"p.gid\"","CursorExpression: missing","expected constant string"},
			{"CursorExpression: \"p.gid\"","CursorExpression: \"p.other\"","no sortable field matches"},
		}{
			root:=searchFixture(searchTestImpl,strings.Replace(searchTestConfig,test.old,test.new,1))
			_,err:=helpers.Discover(root,"grovehub","demo")
			Expect(err).To(MatchError(ContainSubstring(test.reason)),test.new)
		}
		root:=searchFixture(strings.Replace(searchTestImpl,"alias=p","alias=p\n// @Sortable fields=other",1),searchTestConfig)
		_,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(MatchError(ContainSubstring("@Sortable annotation must exactly match")))
		root=searchFixture(strings.Replace(searchTestImpl,"alias=p","alias=p\n// @Filterable fields=gid\n// @Sortable fields=gid\n// @Cursor field=gid",1),searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred());Expect(specs).To(HaveLen(1))
	})
	It("validates implementation ownership and rejects an existing generated method",func(){
		root:=searchFixture(searchTestImpl,searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred())
		dir:=filepath.Join(root,"impl")
		Expect(helpers.ValidateImplementation(dir,specs)).To(Succeed())
		fixture.Write(root,"impl/person.go",strings.Replace(searchTestImpl,"db *gorm.DB","other *gorm.DB",1))
		Expect(helpers.ValidateImplementation(dir,specs)).To(MatchError(ContainSubstring("needs a db field")))
		fixture.Write(root,"impl/person.go",searchTestImpl+"\nfunc (repository *PersonRepository) searchGenerated() {}\n")
		Expect(helpers.ValidateImplementation(dir,specs)).To(MatchError(ContainSubstring("already implements searchGenerated")))
		fixture.Write(root,"impl/person.go","package impl\nfunc broken(")
		Expect(helpers.ValidateImplementation(dir,specs)).NotTo(Succeed())
		Expect(os.Remove(filepath.Join(dir,"person.go"))).To(Succeed())
		Expect(helpers.ValidateImplementation(dir,specs)).NotTo(Succeed())
	})
})
