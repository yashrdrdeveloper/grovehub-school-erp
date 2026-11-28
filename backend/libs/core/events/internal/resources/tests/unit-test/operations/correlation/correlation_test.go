package correlation_test

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"grovehub/libs/core/events/internal/src/operations/correlation"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithID", func() {
	It("propagates correlation without changing the parent context", func() {
		parent := context.Background()
		Expect(correlation.Resolve(correlation.WithID(parent, "request"), "")).To(Equal("request"))
		Expect(correlation.Resolve(parent, "")).To(BeEmpty())
	})
})

var _ = Describe("Resolve", func() {
	It("uses explicit metadata before context and trace", func() {
		ctx := correlation.WithID(context.Background(), "context")
		Expect(correlation.Resolve(ctx, "metadata")).To(Equal("metadata"))
	})

	It("falls back to the trace ID when context is absent or empty", func() {
		traceID, err := trace.TraceIDFromHex("0123456789abcdef0123456789abcdef")
		Expect(err).NotTo(HaveOccurred())
		spanID, err := trace.SpanIDFromHex("0123456789abcdef")
		Expect(err).NotTo(HaveOccurred())
		ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled,
		}))
		Expect(correlation.Resolve(ctx, "")).To(Equal(traceID.String()))
		Expect(correlation.Resolve(correlation.WithID(ctx, ""), "")).To(Equal(traceID.String()))
		Expect(correlation.Resolve(correlation.WithID(ctx, "context"), "")).To(Equal("context"))
	})
})
