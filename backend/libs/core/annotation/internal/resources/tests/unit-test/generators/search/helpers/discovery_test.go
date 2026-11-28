package helpers_test

import (
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("Discover",func(){
	It("rejects annotations on a non-Search method",func(){
		root:=searchFixture(strings.Replace(searchTestImpl,"Search(ctx","Find(ctx",1),searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(MatchError(ContainSubstring("requires a concrete Search method")))
		Expect(specs).To(BeNil())
	})
})
