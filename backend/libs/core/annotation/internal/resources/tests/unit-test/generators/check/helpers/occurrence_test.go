package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/check/helpers"
	contract "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
)

var _ = Describe("ValidateOccurrence",func(){
	It("rejects unknown directives and mismatched targets",func(){
		Expect(helpers.ValidateOccurrence(contract.Occurrence{Directive:contract.Directive{Name:"Unknown"},Target:contract.TypeTarget})).To(MatchError(ContainSubstring("unknown annotation")))
		Expect(helpers.ValidateOccurrence(contract.Occurrence{Directive:contract.Directive{Name:"Entity"},Target:contract.FieldTarget})).To(MatchError(ContainSubstring("cannot annotate target")))
		Expect(helpers.ValidateOccurrence(contract.Occurrence{Directive:contract.Directive{Name:"Component"},Target:contract.TypeTarget})).To(MatchError(ContainSubstring("catalog metadata only")))
	})
})
