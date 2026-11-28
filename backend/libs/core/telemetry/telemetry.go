package telemetry

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"grovehub/libs/core/telemetry/internal/src/contracts"
	"grovehub/libs/core/telemetry/internal/src/operations/runtime"
	"grovehub/libs/core/telemetry/internal/src/helpers"
	"grovehub/libs/core/telemetry/internal/src/operations/tracing"
	"grovehub/libs/core/telemetry/internal/src/config"
	"grovehub/libs/core/telemetry/internal/src/constants"
)

const ServiceNamespace = constants.ServiceNamespace

type Shutdown = contracts.Shutdown
type Service = contracts.Service
type Runtime = runtime.Runtime

var (
	ServiceLauncherEdge = contracts.ServiceLauncherEdge
	ServiceFoundation   = contracts.ServiceFoundation
	ServiceIdentity     = contracts.ServiceIdentity
	ServiceAudit        = contracts.ServiceAudit
)

func ResolveService(value string) Service {
	return helpers.ResolveService(value)
}

func Initialize(ctx context.Context, logger *slog.Logger) Shutdown {
	return runtime.Initialize(ctx, logger)
}

func TracerProvider(serviceName string) trace.TracerProvider {
	return runtime.TracerProvider(serviceName)
}

func MeterProvider(serviceName string) metric.MeterProvider {
	return runtime.MeterProvider(serviceName)
}

func RecordOperation(
	ctx context.Context,
	serviceName string,
	operation string,
	duration time.Duration,
	statusCode int,
) {
	runtime.RecordOperation(ctx, serviceName, operation, duration, statusCode)
}

func RecordHTTPServerRequest(
	ctx context.Context,
	serviceName string,
	method string,
	route string,
	duration time.Duration,
	statusCode int,
) {
	runtime.RecordHTTPServerRequest(ctx, serviceName, method, route, duration, statusCode)
}

func AddActiveHTTPServerRequest(
	ctx context.Context,
	serviceName string,
	method string,
	delta int64,
) {
	runtime.AddActiveHTTPServerRequest(ctx, serviceName, method, delta)
}

func StartSpan(
	ctx context.Context,
	component string,
	operation string,
	options ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return tracing.StartSpan(ctx, component, operation, options...)
}

func ContextWithService(ctx context.Context, service Service) context.Context {
	return tracing.ContextWithService(ctx, service)
}

func ServiceFromContext(ctx context.Context) (Service, bool) {
	return tracing.ServiceFromContext(ctx)
}

func TraceID(ctx context.Context) string {
	return tracing.TraceID(ctx)
}

func SpanID(ctx context.Context) string {
	return tracing.SpanID(ctx)
}

func LogFields(ctx context.Context) []any {
	return tracing.LogFields(ctx)
}

func Logger(logger *slog.Logger, ctx context.Context) *slog.Logger {
	return tracing.Logger(logger, ctx)
}

const (
	DefaultTraceSampleRatio = constants.DefaultTraceSampleRatio
	DefaultMetricExportInterval = constants.DefaultMetricExportInterval
)

func TraceSampleRatio(value string) (float64, error) {
	return config.TraceSampleRatio(value)
}

func MetricExportInterval(value string) (time.Duration, error) {
	return config.MetricExportInterval(value)
}
