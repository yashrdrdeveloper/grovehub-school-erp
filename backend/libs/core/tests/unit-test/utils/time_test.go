package utils_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreutils "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils"
)

var _ = Describe("RFC3339 UTC normalization", func() {
	It("converts offsets to UTC and retains fractional seconds", func() {
		value, err := coreutils.NormalizeRFC3339UTC(" 2026-11-11T15:00:00.123456789+05:30 ")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("2026-11-11T09:30:00.123456789Z"))
	})

	It("rejects malformed timestamps with a domain-neutral error", func() {
		value, err := coreutils.NormalizeRFC3339UTC("not-a-timestamp")
		Expect(value).To(BeEmpty())
		Expect(errors.Is(err, coreutils.ErrInvalidTime)).To(BeTrue())
	})

	It("preserves valid UTC timestamps without adding unnecessary fractions", func() {
		value, err := coreutils.NormalizeRFC3339UTC("2026-11-11T09:30:00Z")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("2026-11-11T09:30:00Z"))
	})
})
