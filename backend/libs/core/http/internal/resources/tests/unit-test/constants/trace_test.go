package constants_test

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/internal/src/constants"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TraceID", func() {
	It("preserves the existing HTTP header name", func() {
		Expect(constants.TraceID).To(Equal("X-Trace-ID"))
	})
})
