package contract_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	contract "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
)

var _ = Describe("Target", func() {
	It("preserves the declared metadata contract", func() {
		Expect([]contract.Target{contract.TypeTarget, contract.InterfaceTarget, contract.MethodTarget, contract.FunctionTarget, contract.FieldTarget, contract.PackageTarget, contract.FileTarget, contract.ParameterTarget}).To(Equal([]contract.Target{"type", "interface", "method", "function", "field", "package", "file", "parameter"}))
	})
})
