package telemetry

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func StartSpan(
	ctx context.Context,
	component string,
	operation string,
	options ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	baseOptions := []trace.SpanStartOption{
		trace.WithAttributes(
			attribute.String("grovehub.component", component),
			attribute.String("grovehub.operation", operation),
		),
	}

	return otel.Tracer("grovehub/" + component).Start(
		ctx,
		component+" "+operation,
		append(baseOptions, options...)...,
	)
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
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return nil
	}

	return []any{
		"traceId", spanContext.TraceID().String(),
		"spanId", spanContext.SpanID().String(),
		"traceFlags", uint8(spanContext.TraceFlags()),
	}
}

func Logger(logger *slog.Logger, ctx context.Context) *slog.Logger {
	fields := LogFields(ctx)
	if len(fields) == 0 {
		return logger
	}

	return logger.With(fields...)
}
