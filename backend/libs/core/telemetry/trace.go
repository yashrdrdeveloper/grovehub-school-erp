package telemetry

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type serviceContextKey struct{}

func StartSpan(
	ctx context.Context,
	component string,
	operation string,
	options ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	service := ResolveService(component)
	baseOptions := []trace.SpanStartOption{
		trace.WithAttributes(
			attribute.String("grovehub.service.key", service.Key),
			attribute.String("grovehub.service.display_name", service.Name),
			attribute.String("grovehub.operation", operation),
		),
	}

	tracer := TracerProvider(service.Key).Tracer("grovehub/" + service.Key)
	spanContext, span := tracer.Start(
		ctx,
		service.Name+" "+operation,
		append(baseOptions, options...)...,
	)

	return ContextWithService(spanContext, service), span
}

func ContextWithService(ctx context.Context, service Service) context.Context {
	return context.WithValue(ctx, serviceContextKey{}, service)
}

func ServiceFromContext(ctx context.Context) (Service, bool) {
	service, ok := ctx.Value(serviceContextKey{}).(Service)
	return service, ok
}

func TraceID(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return ""
	}

	return spanContext.TraceID().String()
}

func SpanID(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return ""
	}

	return spanContext.SpanID().String()
}

func LogFields(ctx context.Context) []any {
	fields := make([]any, 0, 10)
	if service, ok := ServiceFromContext(ctx); ok {
		fields = append(
			fields,
			"service", service.Name,
			"serviceKey", service.Key,
		)
	}

	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		fields = append(
			fields,
			"traceId", spanContext.TraceID().String(),
			"spanId", spanContext.SpanID().String(),
			"traceFlags", uint8(spanContext.TraceFlags()),
		)
	}

	return fields
}

func Logger(logger *slog.Logger, ctx context.Context) *slog.Logger {
	fields := LogFields(ctx)
	if len(fields) == 0 {
		return logger
	}

	return logger.With(fields...)
}
