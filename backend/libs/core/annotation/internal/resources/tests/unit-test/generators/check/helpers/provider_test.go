package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/check/helpers"
	contract "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
)

var _ = Describe("ParseProviderDirective",func(){
	DescribeTable("rejects malformed provider declarations",func(raw string){_,_,err:=helpers.ParseProviderDirective(raw);Expect(err).To(HaveOccurred())},
		Entry("missing category","id=GOOGLE"),Entry("lowercase id","id=google category=mail"),Entry("uppercase category","id=GOOGLE category=MAIL"),Entry("duplicate","id=GOOGLE id=OTHER"),Entry("illegal id","id=GOO-GLE category=mail"),Entry("unknown key","id=GOOGLE kind=mail"))
	It("accepts provider identifiers",func(){id,category,err:=helpers.ParseProviderDirective("id=GOOGLE_2 category=mail-oauth");Expect(err).NotTo(HaveOccurred());Expect(id).To(Equal("GOOGLE_2"));Expect(category).To(Equal("mail-oauth"))})
})
var _ = Describe("Add",func(){
	It("rejects duplicate provider identifiers across files",func(){
		index:=make(helpers.ProviderIndex)
		entry:=contract.Occurrence{Directive:contract.Directive{Name:"Provider",Arguments:"id=GOOGLE category=mail"},File:"first.go",Line:7}
		Expect(index.Add([]contract.Occurrence{entry})).To(Succeed())
		entry.File="second.go"
		Expect(index.Add([]contract.Occurrence{entry})).To(MatchError(ContainSubstring("first.go:7")))
	})
})
