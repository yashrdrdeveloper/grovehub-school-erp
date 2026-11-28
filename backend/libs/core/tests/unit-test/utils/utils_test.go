package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/utils"
)

var _ = Describe("OptionalString", func() {
	It("exposes the empty-to-nil conversion without changing whitespace", func() {
		Expect(utils.OptionalString("")).To(BeNil())
		value := utils.OptionalString("  value  ")
		Expect(value).NotTo(BeNil())
		Expect(*value).To(Equal("  value  "))
	})
})

var _ = Describe("NullableString", func() {
	It("exposes the nil-to-empty conversion and preserves nonnil values", func() {
		Expect(utils.NullableString(nil)).To(BeEmpty())
		value := "  value  "
		Expect(utils.NullableString(&value)).To(Equal(value))
		Expect(utils.NullableString(utils.OptionalString(value))).To(Equal(value))
	})
})
