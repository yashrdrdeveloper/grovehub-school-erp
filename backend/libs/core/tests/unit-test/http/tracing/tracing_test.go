package tracing_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	httptracing "grovehub/libs/core/http/tracing"
	coretelemetry "grovehub/libs/core/telemetry"
)

var _ = Describe("HTTP trace propagation", func() {
	It("creates a child service span on the existing trace", func() {
		request, parentContext, parentSpan := newTracedInboundRequest(http.MethodGet, identityHealthPath)
		defer parentSpan.End()

		var childTraceID string
		var childSpanID string
		handler := httptracing.ServiceHandler(
			identityService,
			identityOperation,
			http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				childTraceID = coretelemetry.TraceID(request.Context())
				childSpanID = coretelemetry.SpanID(request.Context())
				writer.WriteHeader(http.StatusNoContent)
			}),
		)

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusNoContent))
		Expect(childTraceID).To(Equal(coretelemetry.TraceID(parentContext)))
		Expect(childSpanID).NotTo(Equal(coretelemetry.SpanID(parentContext)))
	})

	It("forwards X-Trace-ID and W3C traceparent on outbound HTTP", func() {
		request, parentContext, parentSpan := newTracedOutboundRequest(http.MethodGet, internalServiceURL)
		defer parentSpan.End()

		capturedTraceID, capturedTraceParent := captureOutboundTraceHeaders(request)

		Expect(capturedTraceID).To(Equal(coretelemetry.TraceID(parentContext)))
		Expect(capturedTraceParent).To(ContainSubstring(coretelemetry.TraceID(parentContext)))
	})

	It("creates and forwards a trace when an outbound HTTP call has no parent trace", func() {
		capturedTraceID, capturedTraceParent := captureOutboundTraceHeaders(
			newOutboundRequest(http.MethodGet, backgroundServiceURL),
		)

		Expect(capturedTraceID).NotTo(BeEmpty())
		Expect(capturedTraceParent).NotTo(BeEmpty())
		Expect(capturedTraceParent).To(ContainSubstring(capturedTraceID))
	})
})
