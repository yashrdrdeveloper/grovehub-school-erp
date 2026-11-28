package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover",func(){
	It("retains repository, entity and config identity for generated search methods",func(){
		root:=searchFixture(searchTestImpl,searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred())
		Expect(specs).To(HaveLen(1))
		out:=filepath.Join(root,"output.txt")
		Expect(helpers.Write(out,"grovehub",specs)).To(Succeed())
		content:=fixture.Read(out)
		Expect(content).To(ContainSubstring("func (repository *PersonRepository) searchGenerated"))
		Expect(content).To(ContainSubstring("entity.Person"))
		Expect(content).To(ContainSubstring("search.PersonConfig"))
	})
})
