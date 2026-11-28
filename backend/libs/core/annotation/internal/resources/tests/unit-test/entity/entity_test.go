package entity_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/entity"
)

type record struct { Name string }

var _ = Describe("New", func() {
	It("starts at the zero value", func() {
		Expect(entity.New[record]().Build()).To(Equal(record{}))
	})
})

var _ = Describe("From", func() {
	It("copies the supplied value and supports independent builder branches", func() {
		original := entity.From(record{Name:"original"})
		changed := original.With(func(value *record){value.Name="changed"})
		Expect(original.Build().Name).To(Equal("original"))
		Expect(changed.Build().Name).To(Equal("changed"))
		Expect(changed.With(nil).Build()).To(Equal(changed.Build()))
	})
})
