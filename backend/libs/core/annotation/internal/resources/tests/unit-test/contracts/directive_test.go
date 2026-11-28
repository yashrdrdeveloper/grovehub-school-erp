package contract_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	contract "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
)

var _ = Describe("Directive", func() {
	It("preserves the declared metadata contract", func() {
		directive := contract.Directive{Name: "Entity", Arguments: "school.persons mode=table"}
		Expect(directive.Name).To(Equal("Entity"))
		Expect(directive.Arguments).To(Equal("school.persons mode=table"))
	})
})
