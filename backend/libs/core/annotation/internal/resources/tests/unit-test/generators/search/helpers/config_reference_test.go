package helpers_test

import (
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("Discover",func(){
	It("rejects nonexistent config package",func(){
		root:=searchFixture(strings.Replace(searchTestImpl,"search.PersonConfig","missing.PersonConfig",1),searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).To(MatchError(ContainSubstring("config package")))
		Expect(specs).To(BeNil())
	})
})
