package helpers_test

import (
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("Discover",func(){
	It("rejects declared filter fields outside the configured allowlist",func(){
		root:=searchFixture(strings.Replace(searchTestImpl,"// @Searchable config=search.PersonConfig alias=p","// @Searchable config=search.PersonConfig alias=p\n// @Filterable fields=unknown",1),searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(MatchError(ContainSubstring("must exactly match configured allowlist")))
		Expect(specs).To(BeNil())
	})
})
