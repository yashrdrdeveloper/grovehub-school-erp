package constants_test

import (
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/utils/internal/src/constants"
)

var _ = Describe("Normalization constants", func() {
	It("retains the email limit and international phone pattern", func() {
		Expect(constants.MaxEmailLength).To(Equal(254))
		pattern := regexp.MustCompile(constants.PhonePattern)
		Expect(pattern.MatchString("+919876543210")).To(BeTrue())
		Expect(pattern.MatchString("+0123456789")).To(BeFalse())
	})
})
