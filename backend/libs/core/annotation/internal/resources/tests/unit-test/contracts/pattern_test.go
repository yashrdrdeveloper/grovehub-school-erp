package contract_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	contract "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
)

var _ = Describe("Pattern", func() {
	It("preserves the declared metadata contract", func() {
		pattern := contract.Pattern{Name: "Strategy", Family: "Behavioral", Pattern: "Strategy", Purpose: "selection", Decision: "registry"}
		Expect(pattern).To(Equal(contract.Pattern{"Strategy", "Behavioral", "Strategy", "selection", "registry"}))
	})
})
