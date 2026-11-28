package patterns_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/catalog/patterns"
)

var _ = Describe("All", func() {
	It("includes every creational pattern and its automation decision", func() {
		var names []string

		for _, entry := range patterns.All() {
			if entry.Family != "Creational" {
				continue
			}

			Expect(entry.Pattern).NotTo(BeEmpty())
			Expect(entry.Purpose).NotTo(BeEmpty())
			Expect(entry.Decision).NotTo(BeEmpty())
			names = append(names, entry.Name)
		}

		Expect(names).To(ConsistOf("AbstractFactory", "Builder", "FactoryMethod", "Prototype", "Singleton"))
	})
})
