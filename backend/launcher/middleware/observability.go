package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	httpheaders "grovehub/libs/core/http/headers"
	httplogging "grovehub/libs/core/http/logging"
	coretelemetry "grovehub/libs/core/telemetry"
)

type statusRecorder struct {
	http.ResponseWriter
	status      int
	bytes       int64
	wroteHeader bool
}

func (recorder *statusRecorder) WriteHeader(status int) {
	if recorder.wroteHeader {
		return
	}

	recorder.status = status
	recorder.wroteHeader = true
	recorder.ResponseWriter.WriteHeader(status)
}

func (recorder *statusRecorder) Write(payload []byte) (int, error) {
	if !recorder.wroteHeader {
		recorder.WriteHeader(http.StatusOK)
	}

	written, err := recorder.ResponseWriter.Write(payload)
	recorder.bytes += int64(written)
	return written, err
}

func (recorder *statusRecorder) Unwrap() http.ResponseWriter {
	return recorder.ResponseWriter
}

func RequestObservability(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()
		ctx := request.Context()
		traceID := coretelemetry.TraceID(ctx)
		if traceID != "" {
			writer.Header().Set(httpheaders.TraceID, traceID)
		}

		recorder := &statusRecorder{
			ResponseWriter: writer,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, request)

		duration := time.Since(startedAt)
		span := trace.SpanFromContext(ctx)
		route := strings.TrimPrefix(request.Pattern, request.Method+" ")
		if route == "" {
			route = "unmatched"
		}
		span.SetName("edge " + request.Method + " " + route)
		span.SetAttributes(
			attribute.String("grovehub.component", "edge"),
			attribute.String("http.request.method", request.Method),
			attribute.String("url.path", request.URL.Path),
			attribute.String("http.route", route),
			attribute.Int64("http.request.body.size", request.ContentLength),
			attribute.Int("http.response.status_code", recorder.status),
			attribute.Int64("http.response.body.size", recorder.bytes),
			attribute.Int64("http.server.duration_ms", duration.Milliseconds()),
		)

		if recorder.status >= http.StatusBadRequest {
			span.SetStatus(codes.Error, http.StatusText(recorder.status))
			span.AddEvent(
				"http.request.failed",
				trace.WithAttributes(
					attribute.String("http.request.method", request.Method),
					attribute.String("http.route", route),
					attribute.String("url.path", request.URL.Path),
					attribute.Int("http.response.status_code", recorder.status),
					attribute.Int64("http.response.body.size", recorder.bytes),
				),
			)
		}

		coretelemetry.Logger(logger, ctx).Log(
			ctx,
			httplogging.LevelForStatus(recorder.status),
			"request completed",
			"service", "edge",
			slog.Group(
				"request",
				"method", request.Method,
				"route", route,
				"path", request.URL.Path,
				"contentLength", request.ContentLength,
			),
			slog.Group(
				"response",
				"status", recorder.status,
				"bytes", recorder.bytes,
			),
			"durationMs", duration.Milliseconds(),
		)
	})
}
