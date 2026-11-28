package contract_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	contract "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
)

var _ = Describe("Mode", func() {
	It("preserves the declared metadata contract", func() {
		Expect([]contract.Mode{contract.Generate, contract.Register, contract.Check, contract.Wrap, contract.Tag}).To(Equal([]contract.Mode{"GEN", "REG", "CHECK", "WRAP", "TAG"}))
	})
})
