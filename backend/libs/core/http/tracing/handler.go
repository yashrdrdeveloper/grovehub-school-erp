package tracing

import (
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	coretelemetry "grovehub/libs/core/telemetry"
)

type statusWriter struct {
	http.ResponseWriter
	status      int
	bytes       int64
	wroteHeader bool
}

func (writer *statusWriter) WriteHeader(status int) {
	if writer.wroteHeader {
		return
	}

	writer.status = status
	writer.wroteHeader = true
	writer.ResponseWriter.WriteHeader(status)
}

func (writer *statusWriter) Write(payload []byte) (int, error) {
	if !writer.wroteHeader {
		writer.WriteHeader(http.StatusOK)
	}

	written, err := writer.ResponseWriter.Write(payload)
	writer.bytes += int64(written)
	return written, err
}

func (writer *statusWriter) Unwrap() http.ResponseWriter {
	return writer.ResponseWriter
}

func ServiceHandler(service string, operation string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()
		ctx, span := coretelemetry.StartSpan(
			request.Context(),
			service,
			operation,
			trace.WithSpanKind(trace.SpanKindInternal),
		)
		defer span.End()

		route := strings.TrimPrefix(request.Pattern, request.Method+" ")
		if route == "" {
			route = request.URL.Path
		}

		span.SetAttributes(
			attribute.String("http.request.method", request.Method),
			attribute.String("url.path", request.URL.Path),
			attribute.String("http.route", route),
		)

		recorder := &statusWriter{
			ResponseWriter: writer,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, request.WithContext(ctx))

		span.SetAttributes(
			attribute.Int("http.response.status_code", recorder.status),
			attribute.Int64("http.response.body.size", recorder.bytes),
		)
		if recorder.status >= http.StatusBadRequest {
			span.SetStatus(codes.Error, http.StatusText(recorder.status))
			span.AddEvent(
				"service.operation.failed",
				trace.WithAttributes(
					attribute.String("grovehub.operation", operation),
					attribute.Int("http.response.status_code", recorder.status),
				),
			)
		}

		coretelemetry.RecordOperation(
			ctx,
			service,
			operation,
			time.Since(startedAt),
			recorder.status,
		)
	})
}
