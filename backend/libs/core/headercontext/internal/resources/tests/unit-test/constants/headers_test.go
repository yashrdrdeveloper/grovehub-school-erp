package constants_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/constants"
)

var _ = Describe("Header names", func() {
	It("preserves the inbound metadata wire names", func() {
		Expect(constants.HeaderIdempotencyKey).To(Equal("Idempotency-Key"))
		Expect(constants.HeaderCorrelationID).To(Equal("X-Correlation-ID"))
		Expect(constants.HeaderRequestID).To(Equal("X-Request-ID"))
	})
})
