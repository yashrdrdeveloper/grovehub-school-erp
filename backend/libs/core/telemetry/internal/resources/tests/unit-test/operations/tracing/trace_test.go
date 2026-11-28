package tracing_test

import (
	"bytes"
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"grovehub/libs/core/telemetry/internal/src/contracts"
	"grovehub/libs/core/telemetry/internal/src/operations/tracing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Trace operations",func(){
	It("returns empty identifiers and the original logger without context",func(){
		ctx:=context.Background()
		logger:=slog.New(slog.NewTextHandler(GinkgoWriter,nil))
		Expect(tracing.TraceID(ctx)).To(BeEmpty())
		Expect(tracing.SpanID(ctx)).To(BeEmpty())
		Expect(tracing.LogFields(ctx)).To(BeEmpty())
		_,ok:=tracing.ServiceFromContext(ctx)
		Expect(ok).To(BeFalse())
		Expect(tracing.Logger(logger,ctx)).To(BeIdenticalTo(logger))
		ctx=tracing.ContextWithService(ctx,contracts.ServiceAudit)
		Expect(tracing.LogFields(ctx)).To(Equal([]any{"service","AUDIT","serviceKey","audit"}))
	})
	It("propagates trace context and enriches structured logs",func(){
		previous:=otel.GetTracerProvider()
		provider:=sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
		otel.SetTracerProvider(provider)
		DeferCleanup(func(){otel.SetTracerProvider(previous);Expect(provider.Shutdown(context.Background())).To(Succeed())})
		ctx,parent:=tracing.StartSpan(context.Background(),"edge","parent")
		defer parent.End()
		childCtx,child:=tracing.StartSpan(ctx,"identity","child")
		defer child.End()
		Expect(tracing.TraceID(childCtx)).NotTo(BeEmpty())
		Expect(tracing.TraceID(childCtx)).To(Equal(tracing.TraceID(ctx)))
		Expect(tracing.SpanID(childCtx)).NotTo(Equal(tracing.SpanID(ctx)))
		service,ok:=tracing.ServiceFromContext(childCtx)
		Expect(ok).To(BeTrue())
		Expect(service).To(Equal(contracts.ServiceIdentity))
		var output bytes.Buffer
		logger:=slog.New(slog.NewJSONHandler(&output,nil))
		tracing.Logger(logger,childCtx).Info("correlated")
		Expect(output.String()).To(ContainSubstring(tracing.TraceID(childCtx)))
		Expect(output.String()).To(ContainSubstring(tracing.SpanID(childCtx)))
	})
})
