package entity_test

import (
	coreentity "grovehub/libs/core/annotation"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type sample struct {
	Name  string
	Count int
}

var _ = Describe("Builder", func() {
	It("constructs a typed entity directly", func() {
		actual := coreentity.NewEntity[sample]().With(func(v *sample) { v.Name = "A" }).With(func(v *sample) { v.Count = 2 }).Build()
		Expect(actual).To(Equal(sample{Name: "A", Count: 2}))
	})
	It("copies value structs without changing their source", func() {
		original := sample{Name: "old"}
		updated := coreentity.FromEntity(original).With(func(v *sample) { v.Name = "new" }).Build()
		Expect(original.Name).To(Equal("old"))
		Expect(updated.Name).To(Equal("new"))
	})
	It("accepts a nil change", func() {
		Expect(coreentity.NewEntity[sample]().With(nil).Build()).To(Equal(sample{}))
	})
})
