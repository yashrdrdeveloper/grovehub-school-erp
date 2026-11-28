package helpers_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Optional search directives and static map validation",func(){
	It("rejects present but empty annotations and incomplete or repeated options",func(){
		for _,directive:=range []string{"@Filterable","@Sortable","@Filterable fields","@Filterable fields=","@Sortable fields=gid fields=other"}{
			root:=searchFixture(strings.Replace(searchTestImpl,"alias=p","alias=p\n// "+directive,1),searchTestConfig)
			_,err:=helpers.Discover(root,"grovehub","demo")
			Expect(err).To(HaveOccurred(),directive)
		}
	})
	It("rejects an unkeyed map and requires the explicitly declared cursor's sort expression",func(){
		root:=searchFixture(searchTestImpl,strings.Replace(searchTestConfig,"map[string]coresearch.ScalarField{\"gid\": {Expression: \"p.gid\"}}","map[string]coresearch.ScalarField{{Expression: \"p.gid\"}}",1))
		_,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(MatchError(ContainSubstring("map requires keyed entries")))
		config:=strings.Replace(searchTestConfig,"SortFields: map[string]coresearch.SortField{\"gid\": {Expression: \"p.gid\"}}","SortFields: map[string]coresearch.SortField{\"other\": {Expression: \"p.other\"},\"gid\":{Expression: \"p.gid\"}}",1)
		root=searchFixture(strings.Replace(searchTestImpl,"alias=p","alias=p\n// @Cursor field=gid",1),config)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred());Expect(specs).To(HaveLen(1))
		Expect(helpers.Write(fixture.Write(root,"out.go","package impl"),"grovehub",nil)).To(MatchError("no search specifications to generate"))
	})
})
