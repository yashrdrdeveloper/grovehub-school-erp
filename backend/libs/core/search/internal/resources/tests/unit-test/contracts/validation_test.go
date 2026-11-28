package contracts_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/search/internal/src/contracts"
)

var _ = Describe("FieldSpec",func(){
	It("allows resource-owned normalization and operator policy",func(){
		spec:=contracts.FieldSpec{Operators:map[contracts.Operator]struct{}{contracts.OperatorEquals:{}},Normalizer:func(_ contracts.Operator,value string)(string,error){return value,nil}}
		Expect(spec.Operators).To(HaveKey(contracts.OperatorEquals))
		value,err:=spec.Normalizer(contracts.OperatorEquals,"value")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("value"))
	})
})
