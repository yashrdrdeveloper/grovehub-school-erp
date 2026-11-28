package contract_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	contract "grovehub/libs/core/annotation/internal/src/contracts"
)

var _ = Describe("Definition", func() {
	It("preserves the declared metadata contract", func() {
		definition := contract.Definition{Name: "Retry", UseCases: []contract.UseCase{{Name: "Retry", Target: contract.MethodTarget}, {Name: "Retry", Target: contract.FunctionTarget}}}
		Expect(definition.UseCases).To(HaveLen(2))
		Expect(definition.UseCases[0].Target).NotTo(Equal(definition.UseCases[1].Target))
	})
})
