package telemetry_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coretelemetry "grovehub/libs/core/telemetry"
)

var _ = Describe("Trace context", func() {
	It("keeps one trace ID across logical service spans", func() {
		rootContext, rootSpan := coretelemetry.StartSpan(
			context.Background(),
			coretelemetry.ServiceLauncherEdge.Key,
			"GET /identity/health",
		)
		defer rootSpan.End()

		childContext, childSpan := coretelemetry.StartSpan(rootContext, testComponent, testOperation)
		defer childSpan.End()

		Expect(coretelemetry.TraceID(rootContext)).NotTo(BeEmpty())
		Expect(coretelemetry.TraceID(childContext)).To(Equal(coretelemetry.TraceID(rootContext)))
		Expect(coretelemetry.SpanID(childContext)).NotTo(Equal(coretelemetry.SpanID(rootContext)))

		rootService, rootServiceFound := coretelemetry.ServiceFromContext(rootContext)
		childService, childServiceFound := coretelemetry.ServiceFromContext(childContext)
		Expect(rootServiceFound).To(BeTrue())
		Expect(childServiceFound).To(BeTrue())
		Expect(rootService).To(Equal(coretelemetry.ServiceLauncherEdge))
		Expect(childService).To(Equal(coretelemetry.ServiceIdentity))
	})

	It("adds service trace and span identifiers to structured logs", func() {
		ctx, span := coretelemetry.StartSpan(context.Background(), testComponent, testOperation)
		defer span.End()

		logger, buffer := newJSONLogger()
		coretelemetry.Logger(logger, ctx).Info("trace correlated log")
		payload := decodeLog(buffer)

		Expect(payload["service"]).To(Equal(coretelemetry.ServiceIdentity.Name))
		Expect(payload["serviceKey"]).To(Equal(coretelemetry.ServiceIdentity.Key))
		Expect(payload["traceId"]).To(Equal(coretelemetry.TraceID(ctx)))
		Expect(payload["spanId"]).To(Equal(coretelemetry.SpanID(ctx)))
	})
})
