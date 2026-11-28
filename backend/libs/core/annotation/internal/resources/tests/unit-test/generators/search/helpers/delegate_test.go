package helpers_test

import (
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("Discover",func(){
	It("rejects changed context forwarding",func(){
		root:=searchFixture(strings.Replace(searchTestImpl,"repository.searchGenerated(ctx, query)","repository.searchGenerated(context.Background(), query)",1),searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(MatchError(ContainSubstring("forward ctx and query unchanged")))
		Expect(specs).To(BeNil())
	})
})
