package strings_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/operations/strings"
)

var _ = Describe("OptionalString", func() {
	It("maps only the empty string to nil", func() {
		Expect(strings.OptionalString("")).To(BeNil())
		for _, value := range []string{"value", "  value  ", " ", "é"} {
			got := strings.OptionalString(value)
			Expect(got).NotTo(BeNil())
			Expect(*got).To(Equal(value))
		}
	})

	It("returns independent pointers without shared state", func() {
		first := strings.OptionalString("value")
		second := strings.OptionalString("value")
		*first = "changed"
		Expect(*second).To(Equal("value"))
	})
})

var _ = Describe("NullableString", func() {
	It("maps nil to empty and preserves all nonnil values", func() {
		Expect(strings.NullableString(nil)).To(BeEmpty())
		for _, value := range []string{"", "value", "  value  ", " ", "é"} {
			Expect(strings.NullableString(&value)).To(Equal(value))
		}
	})
})
