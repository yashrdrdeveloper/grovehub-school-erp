package runtime_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/operations/runtime"
)

var sdkFailure=errors.New("SDK failure")

type failedMeter struct{ noop.Meter }
func (failedMeter) Float64Histogram(string,...metric.Float64HistogramOption)(metric.Float64Histogram,error){return nil,sdkFailure}
func (failedMeter) Int64Counter(string,...metric.Int64CounterOption)(metric.Int64Counter,error){return nil,sdkFailure}
func (failedMeter) Int64UpDownCounter(string,...metric.Int64UpDownCounterOption)(metric.Int64UpDownCounter,error){return nil,sdkFailure}
func (failedMeter) Int64ObservableGauge(string,...metric.Int64ObservableGaugeOption)(metric.Int64ObservableGauge,error){return nil,sdkFailure}

type failedTraceExporter struct{}
func (failedTraceExporter) ExportSpans(context.Context,[]sdktrace.ReadOnlySpan)error{return nil}
func (failedTraceExporter) Shutdown(context.Context)error{return sdkFailure}

type failedMetricExporter struct{}
func (failedMetricExporter) Temporality(kind sdkmetric.InstrumentKind)metricdata.Temporality{return sdkmetric.DefaultTemporalitySelector(kind)}
func (failedMetricExporter) Aggregation(kind sdkmetric.InstrumentKind)sdkmetric.Aggregation{return sdkmetric.DefaultAggregationSelector(kind)}
func (failedMetricExporter) Export(context.Context,*metricdata.ResourceMetrics)error{return nil}
func (failedMetricExporter) ForceFlush(context.Context)error{return nil}
func (failedMetricExporter) Shutdown(context.Context)error{return sdkFailure}

var _=Describe("Telemetry SDK boundaries",func(){
	It("degrades safely when SDK initialization fails",func(){
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_ENDPOINT","http://127.0.0.1:4318")
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT","")
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT","")
		var output bytes.Buffer
		logger:=slog.New(slog.NewTextHandler(&output,nil))
		shutdown:=runtime.InitializeWithDependencies(context.Background(),logger,runtime.Dependencies{
			TraceExporter:func(context.Context)(sdktrace.SpanExporter,error){return nil,sdkFailure},
			MetricExporter:func(context.Context)(sdkmetric.Exporter,error){return nil,sdkFailure},
			StartRuntimeMetrics:func(metric.MeterProvider)error{return sdkFailure},
			Meter:func(metric.MeterProvider,string)metric.Meter{return failedMeter{}},
		})
		otel.Handle(sdkFailure)
		Expect(output.String()).To(ContainSubstring("trace exporter initialization failed"))
		Expect(output.String()).To(ContainSubstring("metric exporter initialization failed"))
		Expect(output.String()).To(ContainSubstring("runtime metrics initialization failed"))
		Expect(output.String()).To(ContainSubstring("failed to create service availability metric"))
		Expect(output.String()).To(ContainSubstring("OpenTelemetry export failed"))
		Expect(shutdown(context.Background())).To(Succeed())
	})
	It("reports both provider shutdown failures",func(){
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_ENDPOINT","http://127.0.0.1:4318")
		shutdown:=runtime.InitializeWithDependencies(context.Background(),slog.New(slog.NewTextHandler(GinkgoWriter,nil)),runtime.Dependencies{
			TraceExporter:func(context.Context)(sdktrace.SpanExporter,error){return failedTraceExporter{},nil},
			MetricExporter:func(context.Context)(sdkmetric.Exporter,error){return failedMetricExporter{},nil},
			StartRuntimeMetrics:func(metric.MeterProvider)error{return nil},
			Meter:func(provider metric.MeterProvider,name string)metric.Meter{return provider.Meter(name)},
		})
		ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()
		err:=shutdown(ctx)
		Expect(err).To(MatchError(ContainSubstring("meter provider")))
		Expect(err).To(MatchError(ContainSubstring("tracer provider")))
	})
	It("exports actual traces and metrics to a local OTLP collector",func(){
		var traces,metrics atomic.Int32
		server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
			switch r.URL.Path{case "/v1/traces":traces.Add(1);case "/v1/metrics":metrics.Add(1);default:w.WriteHeader(404);return}
			w.Header().Set("Content-Type","application/x-protobuf")
			w.WriteHeader(200)
		}))
		defer server.Close()
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_ENDPOINT",server.URL)
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT","")
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT","")
		GinkgoT().Setenv("GROVEHUB_TRACE_SAMPLE_RATIO","1")
		GinkgoT().Setenv("GROVEHUB_METRIC_EXPORT_INTERVAL","15s")
		shutdown:=runtime.InitializeWithDependencies(context.Background(),slog.New(slog.NewTextHandler(GinkgoWriter,nil)),runtime.Dependencies{
			TraceExporter:func(ctx context.Context)(sdktrace.SpanExporter,error){return otlptracehttp.New(ctx)},
			MetricExporter:func(ctx context.Context)(sdkmetric.Exporter,error){return otlpmetrichttp.New(ctx)},
			StartRuntimeMetrics:func(metric.MeterProvider)error{return nil},
			Meter:func(provider metric.MeterProvider,name string)metric.Meter{return provider.Meter(name)},
		})
		ctx,span:=runtime.TracerProvider("identity").Tracer("test").Start(context.Background(),"operation")
		runtime.RecordOperation(ctx,"identity","read",time.Millisecond,200)
		span.End()
		deadline,cancel:=context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()
		Expect(shutdown(deadline)).To(Succeed())
		Expect(traces.Load()).To(BeNumerically(">",0))
		Expect(metrics.Load()).To(BeNumerically(">",0))
	})
})
