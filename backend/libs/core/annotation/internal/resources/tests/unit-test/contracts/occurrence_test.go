package contract_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	contract "grovehub/libs/core/annotation/internal/src/contracts"
)

var _ = Describe("Occurrence", func() {
	It("preserves the declared metadata contract", func() {
		occurrence := contract.Occurrence{Directive: contract.Directive{Name: "Entity"}, Target: contract.TypeTarget, File: "person.go", Line: 5, Column: 1}
		Expect(occurrence.Directive.Name).To(Equal("Entity"))
		Expect(occurrence.Target).To(Equal(contract.TypeTarget))
		Expect(occurrence.File).To(Equal("person.go"))
		Expect(occurrence.Line).To(Equal(5))
		Expect(occurrence.Column).To(Equal(1))
	})
})
