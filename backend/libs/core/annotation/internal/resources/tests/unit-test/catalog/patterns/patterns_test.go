package patterns_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/patterns"
)

var _ = Describe("All", func() {
	It("returns all 23 patterns in family order without sharing writable storage", func() {
		first := patterns.All()
		Expect(first).To(HaveLen(23))
		Expect(first[0].Family).To(Equal("Creational"))
		Expect(first[5].Family).To(Equal("Structural"))
		Expect(first[12].Family).To(Equal("Behavioral"))
		Expect(patterns.Entries).To(Equal(first))

		first[0].Name = "changed"

		Expect(patterns.All()[0].Name).To(Equal("AbstractFactory"))
	})
})
