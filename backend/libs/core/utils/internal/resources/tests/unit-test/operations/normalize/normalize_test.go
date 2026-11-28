package normalize_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/operations/normalize"
)

var _ = Describe("NormalizeEmail", func() {
	It("trims and lowercases plain addresses", func() {
		value, err := normalize.NormalizeEmail("  Ada@Example.COM  ")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("ada@example.com"))
	})

	It("rejects empty, oversized, malformed and display-name addresses", func() {
		for _, input := range []string{"", strings.Repeat("a", 255), "not-an-email", "Ada <ada@example.com>"} {
			value, err := normalize.NormalizeEmail(input)
			Expect(value).To(BeEmpty())
			Expect(err).To(BeIdenticalTo(contracts.ErrInvalidEmail))
		}
	})
})

var _ = Describe("NormalizePhone", func() {
	It("normalizes international phone numbers", func() {
		value, err := normalize.NormalizePhone(" +91 (98765) 43210 ")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("+919876543210"))
	})

	It("rejects invalid characters and numbers outside the international pattern", func() {
		for _, input := range []string{"+91abc", "9876543210", "+0123456789", "+123", "+1234567890123456"} {
			value, err := normalize.NormalizePhone(input)
			Expect(value).To(BeEmpty())
			Expect(err).To(BeIdenticalTo(contracts.ErrInvalidPhone))
		}
	})
})

var _ = Describe("NormalizePhoneCharacters", func() {
	It("removes supported separators without requiring an international prefix", func() {
		for input, expected := range map[string]string{" (123) 456-78 ":"12345678", "+1-234":"+1234", "":""} {
			value, err := normalize.NormalizePhoneCharacters(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal(expected))
		}
	})

	It("rejects misplaced plus signs and unsupported characters", func() {
		for _, input := range []string{"12+34", "12.34", "12\t34", "１２"} {
			value, err := normalize.NormalizePhoneCharacters(input)
			Expect(value).To(BeEmpty())
			Expect(err).To(BeIdenticalTo(contracts.ErrInvalidPhone))
		}
	})
})

var _ = Describe("NormalizeUUIDv7", func() {
	It("returns canonical lowercase UUIDv7 strings", func() {
		value, err := normalize.NormalizeUUIDv7(" 01890ABC-1234-7ABC-8DEF-0123456789AB ")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("01890abc-1234-7abc-8def-0123456789ab"))
	})

	It("rejects malformed, nil and other UUID versions", func() {
		for _, input := range []string{"invalid", "00000000-0000-0000-0000-000000000000", "01890abc-1234-4abc-8def-0123456789ab"} {
			value, err := normalize.NormalizeUUIDv7(input)
			Expect(value).To(BeEmpty())
			Expect(err).To(BeIdenticalTo(contracts.ErrInvalidUUIDv7))
		}
	})
})
