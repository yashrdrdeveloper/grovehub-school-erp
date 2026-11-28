package validation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/validation"
)

var _ = Describe("Operator policies",func(){
	It("provides explicit operator allowlists and cardinality checks",func(){
		Expect(validation.EqualityOperators()).To(HaveLen(4))
		Expect(validation.TextOperators()).To(HaveLen(8))
		Expect(validation.TimeOperators()).To(HaveLen(5))
		Expect(validation.IsPartialOperator(contracts.OperatorContains)).To(BeTrue())
		Expect(validation.IsPartialOperator(contracts.OperatorStartsWith)).To(BeTrue())
		Expect(validation.IsPartialOperator(contracts.OperatorEquals)).To(BeFalse())
		for _,item:=range []struct{op contracts.Operator;count int;valid bool}{
			{contracts.OperatorExists,0,true},{contracts.OperatorNotExists,1,false},{contracts.OperatorIn,2,true},{contracts.OperatorNotIn,0,false},{contracts.OperatorEquals,1,true},{contracts.OperatorEquals,2,false},
		}{Expect(validation.ValidValueCount(item.op,item.count)).To(Equal(item.valid))}
	})
})
