package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreutils "grovehub/libs/core/utils"
)

var _ = Describe("Normalization facade", func() {
	It("exposes all renamed normalization operations", func() {
		value, err := coreutils.NormalizeEmail(" ADA@example.com ")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("ada@example.com"))
		value, err = coreutils.NormalizePhone("+91 9876543210")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("+919876543210"))
		value, err = coreutils.NormalizePhoneCharacters("(123) 456")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("123456"))
		value, err = coreutils.NormalizeUUIDv7("01890ABC-1234-7ABC-8DEF-0123456789AB")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("01890abc-1234-7abc-8def-0123456789ab"))
		_, err = coreutils.NormalizeEmail("invalid")
		Expect(err).To(BeIdenticalTo(coreutils.ErrInvalidEmail))
		_, err = coreutils.NormalizePhone("invalid")
		Expect(err).To(BeIdenticalTo(coreutils.ErrInvalidPhone))
		_, err = coreutils.NormalizeUUIDv7("invalid")
		Expect(err).To(BeIdenticalTo(coreutils.ErrInvalidUUIDv7))
	})
})
