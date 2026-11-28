package constants_test

import (
	"grovehub/libs/core/http/internal/src/constants"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TraceID", func() {
	It("preserves the existing HTTP header name", func() {
		Expect(constants.TraceID).To(Equal("X-Trace-ID"))
	})
})
