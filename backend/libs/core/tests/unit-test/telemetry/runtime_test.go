package telemetry_test

import (
	"context"
	"io"
	"log/slog"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coretelemetry "grovehub/libs/core/telemetry"
)

var _ = Describe("Telemetry runtime facade",func(){
	It("initializes providers, records metrics and shuts down",func(){
		for _, key:=range []string{"OTEL_EXPORTER_OTLP_ENDPOINT","OTEL_EXPORTER_OTLP_TRACES_ENDPOINT","OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"} { GinkgoT().Setenv(key,"") }
		ctx:=context.Background()
		shutdown:=coretelemetry.Initialize(ctx,slog.New(slog.NewTextHandler(io.Discard,nil)))
		DeferCleanup(func(){
			deadline,cancel:=context.WithTimeout(ctx,5*time.Second)
			defer cancel()
			Expect(shutdown(deadline)).To(Succeed())
		})
		Expect(coretelemetry.TracerProvider("identity")).NotTo(BeIdenticalTo(coretelemetry.TracerProvider("foundation")))
		Expect(coretelemetry.MeterProvider("identity")).NotTo(BeIdenticalTo(coretelemetry.MeterProvider("foundation")))
		ctx,span:=coretelemetry.StartSpan(ctx,"identity","test")
		defer span.End()
		Expect(coretelemetry.LogFields(ctx)).To(ContainElements("service","IDENTITY","traceId",coretelemetry.TraceID(ctx)))
		coretelemetry.RecordOperation(ctx,"identity","test",time.Millisecond,200)
		coretelemetry.RecordHTTPServerRequest(ctx,"identity","GET","/test",time.Millisecond,500)
		coretelemetry.AddActiveHTTPServerRequest(ctx,"identity","GET",1)
		coretelemetry.AddActiveHTTPServerRequest(ctx,"identity","GET",-1)
	})
})
