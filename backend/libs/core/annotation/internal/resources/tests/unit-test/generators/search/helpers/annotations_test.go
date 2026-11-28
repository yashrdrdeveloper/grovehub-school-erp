package helpers_test

import (
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("Discover",func(){
	It("rejects duplicate searchable options",func(){
		root:=searchFixture(strings.Replace(searchTestImpl,"alias=p","alias=p alias=q",1),searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(MatchError(ContainSubstring("invalid or repeated @Searchable option")))
		Expect(specs).To(BeNil())
	})
})
