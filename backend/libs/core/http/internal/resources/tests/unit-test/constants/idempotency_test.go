package constants_test

import (
	"grovehub/libs/core/http/internal/src/constants"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IdempotencyKey", func() {
	It("preserves the existing HTTP header name", func() {
		Expect(constants.IdempotencyKey).To(Equal("Idempotency-Key"))
	})
})
