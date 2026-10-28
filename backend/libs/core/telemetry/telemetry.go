package telemetry

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"grovehub/libs/core/telemetry/contracts"
	"grovehub/libs/core/telemetry/helpers"
	"grovehub/libs/core/telemetry/mappers"
)

const ServiceNamespace = contracts.ServiceNamespace

type Shutdown = contracts.Shutdown
type Service = contracts.Service
type Runtime = helpers.Runtime

var (
	ServiceLauncherEdge = contracts.ServiceLauncherEdge
	ServiceFoundation   = contracts.ServiceFoundation
	ServiceIdentity     = contracts.ServiceIdentity
	ServiceMailProvider = contracts.ServiceMailProvider
)

func ResolveService(value string) Service {
	return mappers.ResolveService(value)
}

func Initialize(ctx context.Context, logger *slog.Logger) Shutdown {
	return helpers.Initialize(ctx, logger)
}

func TracerProvider(serviceName string) trace.TracerProvider {
	return helpers.TracerProvider(serviceName)
}

func MeterProvider(serviceName string) metric.MeterProvider {
	return helpers.MeterProvider(serviceName)
}

func RecordOperation(
	ctx context.Context,
	serviceName string,
	operation string,
	duration time.Duration,
	statusCode int,
) {
	helpers.RecordOperation(ctx, serviceName, operation, duration, statusCode)
}

func RecordHTTPServerRequest(
	ctx context.Context,
	serviceName string,
	method string,
	route string,
	duration time.Duration,
	statusCode int,
) {
	helpers.RecordHTTPServerRequest(ctx, serviceName, method, route, duration, statusCode)
}

func AddActiveHTTPServerRequest(
	ctx context.Context,
	serviceName string,
	method string,
	delta int64,
) {
	helpers.AddActiveHTTPServerRequest(ctx, serviceName, method, delta)
}

func StartSpan(
	ctx context.Context,
	component string,
	operation string,
	options ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return helpers.StartSpan(ctx, component, operation, options...)
}

func ContextWithService(ctx context.Context, service Service) context.Context {
	return helpers.ContextWithService(ctx, service)
}

func ServiceFromContext(ctx context.Context) (Service, bool) {
	return helpers.ServiceFromContext(ctx)
}

func TraceID(ctx context.Context) string {
	return helpers.TraceID(ctx)
}

func SpanID(ctx context.Context) string {
	return helpers.SpanID(ctx)
}

func LogFields(ctx context.Context) []any {
	return helpers.LogFields(ctx)
}

func Logger(logger *slog.Logger, ctx context.Context) *slog.Logger {
	return helpers.Logger(logger, ctx)
}
