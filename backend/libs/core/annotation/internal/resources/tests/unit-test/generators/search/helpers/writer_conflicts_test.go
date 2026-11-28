package helpers_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Search writer ownership conflicts",func(){
	It("rejects config aliases resolving to different imports",func(){
		root:=searchFixture(searchTestImpl,searchTestConfig)
		first,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred())
		second,err:=helpers.Discover(root,"grovehub","other")
		Expect(err).NotTo(HaveOccurred())
		Expect(helpers.Write(filepath.Join(root,"out.go"),"grovehub",append(first,second...))).To(MatchError(ContainSubstring("conflicting config package alias")))
	})
	It("rejects mixed service ownership even with local configuration variables",func(){
		root:=searchFixture(strings.Replace(searchTestImpl,"config=search.PersonConfig","config=PersonConfig",1),searchTestConfig)
		fixture.Write(root,"impl/config.go",strings.Replace(searchTestConfig,"package search","package impl",1))
		first,err:=helpers.Discover(root,"grovehub","demo");Expect(err).NotTo(HaveOccurred())
		second,err:=helpers.Discover(root,"grovehub","other");Expect(err).NotTo(HaveOccurred())
		Expect(helpers.Write(filepath.Join(root,"out.go"),"grovehub",append(first,second...))).To(MatchError("mixed service annotations"))
	})
})
