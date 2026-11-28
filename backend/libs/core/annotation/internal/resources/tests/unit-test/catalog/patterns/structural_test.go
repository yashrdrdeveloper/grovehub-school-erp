package patterns_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/patterns"
)

var _ = Describe("All", func() {
	It("includes every structural pattern and its automation decision", func() {
		var names []string

		for _, entry := range patterns.All() {
			if entry.Family != "Structural" {
				continue
			}

			Expect(entry.Pattern).NotTo(BeEmpty())
			Expect(entry.Purpose).NotTo(BeEmpty())
			Expect(entry.Decision).NotTo(BeEmpty())
			names = append(names, entry.Name)
		}

		Expect(names).To(ConsistOf("Adapter", "Bridge", "Composite", "Decorator", "Facade", "Flyweight", "Proxy"))
	})
})
