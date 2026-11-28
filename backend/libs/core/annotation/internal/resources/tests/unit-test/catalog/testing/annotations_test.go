package catalogtesting_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	catalog "grovehub/libs/core/annotation/internal/src/catalog/testing"
	contract "grovehub/libs/core/annotation/internal/src/contracts"
	"grovehub/libs/core/annotation/internal/src/registry"
)

var _ = Describe("UseCases", func() {
	It("defines complete, registered testing annotation metadata", func() {
		Expect(catalog.UseCases).NotTo(BeEmpty())

		for _, entry := range catalog.UseCases {
			Expect(entry.Name).NotTo(BeEmpty())
			Expect(entry.Syntax).To(HavePrefix("@" + entry.Name))
			Expect(entry.Category).NotTo(BeEmpty())
			Expect(entry.Purpose).NotTo(BeEmpty())
			Expect(entry.Mode).NotTo(BeEmpty())
			Expect(entry.Roadmap).NotTo(BeEmpty())
			Expect(entry.Target).To(BeElementOf(
				contract.TypeTarget, contract.InterfaceTarget, contract.MethodTarget,
				contract.FunctionTarget, contract.FieldTarget, contract.PackageTarget,
				contract.FileTarget, contract.ParameterTarget,
			))
			Expect(entry.Status).To(BeElementOf(
				contract.Existing, contract.Gombok, contract.Proposed,
				contract.Optional, contract.Discouraged, contract.Pilot,
			))

			definition, found := registry.Lookup(entry.Name)
			Expect(found).To(BeTrue())
			Expect(definition.UseCases).To(ContainElement(entry))
		}
	})
})
