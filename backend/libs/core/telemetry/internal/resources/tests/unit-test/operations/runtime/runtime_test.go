package runtime_test

import (
	"context"
	"io"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/operations/runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Initialize",func(){
	It("creates independent providers and releases them without an external exporter",func(){
		for _,key:=range []string{"OTEL_EXPORTER_OTLP_ENDPOINT","OTEL_EXPORTER_OTLP_TRACES_ENDPOINT","OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"}{GinkgoT().Setenv(key,"")}
		GinkgoT().Setenv("GROVEHUB_TRACE_SAMPLE_RATIO","1")
		GinkgoT().Setenv("GROVEHUB_METRIC_EXPORT_INTERVAL","15s")
		ctx:=context.Background()
		logger:=slog.New(slog.NewTextHandler(io.Discard,nil))
		shutdown:=runtime.Initialize(ctx,logger)
		DeferCleanup(func(){
			deadline,cancel:=context.WithTimeout(ctx,5*time.Second)
			defer cancel()
			Expect(shutdown(deadline)).To(Succeed())
		})
		Expect(runtime.TracerProvider("identity")).NotTo(BeIdenticalTo(runtime.TracerProvider("foundation")))
		Expect(runtime.MeterProvider("audit")).NotTo(BeIdenticalTo(runtime.MeterProvider("identity")))
		Expect(runtime.TracerProvider("unregistered")).To(BeIdenticalTo(otel.GetTracerProvider()))
		Expect(runtime.MeterProvider("unregistered")).To(BeIdenticalTo(otel.GetMeterProvider()))
		_,span:=runtime.TracerProvider("identity").Tracer("test").Start(ctx,"operation")
		Expect(span.SpanContext().IsValid()).To(BeTrue())
		span.End()
		for _,code:=range []int{200,400,500}{
			runtime.RecordOperation(ctx,"identity","test",time.Millisecond,code)
			runtime.RecordHTTPServerRequest(ctx,"identity","GET","/test",time.Millisecond,code)
		}
		runtime.AddActiveHTTPServerRequest(ctx,"identity","GET",1)
		runtime.AddActiveHTTPServerRequest(ctx,"identity","GET",-1)
		runtime.RecordOperation(ctx,"unregistered","test",time.Millisecond,500)
		runtime.RecordHTTPServerRequest(ctx,"unregistered","GET","/test",time.Millisecond,500)
		runtime.AddActiveHTTPServerRequest(ctx,"unregistered","GET",1)
	})
})
