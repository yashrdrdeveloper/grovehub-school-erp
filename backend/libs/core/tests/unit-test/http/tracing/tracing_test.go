package tracing_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/trace"

	httptracing "grovehub/libs/core/http/tracing"
	coretelemetry "grovehub/libs/core/telemetry"
)

var _ = Describe("HTTP trace propagation", func() {
	It("creates launcher client and identity server spans on one trace", func() {
		request, parentContext, parentSpan := newTracedInboundRequest(http.MethodGet, identityHealthPath)
		defer parentSpan.End()

		var childTraceID string
		var childSpanID string
		var childService coretelemetry.Service
		var childServicePresent bool
		handler := httptracing.ServiceHandler(
			identityService,
			identityOperation,
			http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				childTraceID = coretelemetry.TraceID(request.Context())
				childSpanID = coretelemetry.SpanID(request.Context())
				childService, childServicePresent = coretelemetry.ServiceFromContext(request.Context())
				writer.WriteHeader(http.StatusNoContent)
			}),
		)

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusNoContent))
		Expect(childTraceID).To(Equal(coretelemetry.TraceID(parentContext)))
		Expect(childSpanID).NotTo(Equal(coretelemetry.SpanID(parentContext)))
		Expect(childServicePresent).To(BeTrue())
		Expect(childService).To(Equal(coretelemetry.ServiceIdentity))

		endedSpans := testSpanRecorder.Ended()
		var clientSpanID trace.SpanID
		var serverParentSpanID trace.SpanID
		var clientParentSpanID trace.SpanID
		var clientFound bool
		var serverFound bool
		for _, endedSpan := range endedSpans {
			switch endedSpan.SpanKind() {
			case trace.SpanKindClient:
				clientFound = true
				clientSpanID = endedSpan.SpanContext().SpanID()
				clientParentSpanID = endedSpan.Parent().SpanID()
			case trace.SpanKindServer:
				serverFound = true
				serverParentSpanID = endedSpan.Parent().SpanID()
			}
		}

		Expect(clientFound).To(BeTrue())
		Expect(serverFound).To(BeTrue())
		Expect(clientParentSpanID.String()).To(Equal(coretelemetry.SpanID(parentContext)))
		Expect(serverParentSpanID).To(Equal(clientSpanID))
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
