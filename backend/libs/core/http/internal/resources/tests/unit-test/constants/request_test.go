package constants_test

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/internal/src/constants"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestID", func() {
	It("preserves the existing HTTP header name", func() {
		Expect(constants.RequestID).To(Equal("X-Request-ID"))
		Expect(constants.CorrelationID).To(Equal("X-Correlation-ID"))
	})
})
