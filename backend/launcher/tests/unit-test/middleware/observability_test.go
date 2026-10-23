package middleware_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/codes"

	httpheaders "grovehub/libs/core/http/headers"
)

var _ = Describe("Edge request observability", func() {
	It("returns the trace ID and correlates an HTTP failure with an error span", func() {
		observed := serveObservedRequest(http.StatusNotFound)
		endedSpans := spanRecorder.Ended()

		Expect(observed.recorder.Header().Get(httpheaders.TraceID)).To(Equal(observed.traceID))
		Expect(observed.log).To(HaveKeyWithValue("level", "ERROR"))
		Expect(observed.log).To(HaveKeyWithValue("traceId", observed.traceID))
		Expect(observed.log).To(HaveKeyWithValue("spanId", observed.spanID))
		Expect(observed.log["request"]).To(HaveKeyWithValue("path", failedPath))
		Expect(observed.log["response"]).To(HaveKeyWithValue("status", float64(http.StatusNotFound)))
		Expect(endedSpans).To(HaveLen(1))
		Expect(endedSpans[0].Status().Code).To(Equal(codes.Error))
		Expect(endedSpans[0].Events()).NotTo(BeEmpty())
	})
})
