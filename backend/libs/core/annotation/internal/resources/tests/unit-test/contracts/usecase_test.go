package contract_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	contract "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
)

var _ = Describe("UseCase", func() {
	It("preserves the declared metadata contract", func() {
		entry := contract.UseCase{Name: "Retry", Syntax: "@Retry", Category: "external", Purpose: "retry", Mode: "REG/WRAP", Status: contract.Proposed, Roadmap: "P4", Target: contract.MethodTarget}
		Expect(entry.Mode).To(Equal("REG/WRAP"))
		Expect(entry.Status).To(Equal(contract.Proposed))
		Expect(entry.Target).To(Equal(contract.MethodTarget))
	})
})
