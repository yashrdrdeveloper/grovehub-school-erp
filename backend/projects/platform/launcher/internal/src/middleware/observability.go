package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
 "grovehub/launcher/internal/src/middleware/helpers"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"

	coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
)

func RequestObservability(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()
		edgeService := coretelemetry.ServiceLauncherEdge
		ctx := coretelemetry.ContextWithService(request.Context(), edgeService)
		request = request.WithContext(ctx)

		coretelemetry.AddActiveHTTPServerRequest(ctx, edgeService.Key, request.Method, 1)
		defer coretelemetry.AddActiveHTTPServerRequest(ctx, edgeService.Key, request.Method, -1)

		traceID := coretelemetry.TraceID(ctx)
		if traceID != "" {
			writer.Header().Set(corehttp.TraceID, traceID)
		}

		recorder := &helpers.Recorder{
			ResponseWriter: writer,
			Status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, request)

		duration := time.Since(startedAt)
		span := trace.SpanFromContext(ctx)
		route := strings.TrimPrefix(request.Pattern, request.Method+" ")
		if route == "" {
			route = "unmatched"
		}
		if route == "/" && request.URL.Path != "/" && recorder.Status == http.StatusNotFound {
			route = "unmatched"
		}

		span.SetName(edgeService.Name + " " + request.Method + " " + route)
		span.SetAttributes(
			attribute.String("grovehub.service.key", edgeService.Key),
			attribute.String("grovehub.service.display_name", edgeService.Name),
			attribute.String("http.request.method", request.Method),
			attribute.String("url.path", request.URL.Path),
			attribute.String("http.route", route),
			attribute.Int64("http.request.body.size", request.ContentLength),
			attribute.Int("http.response.status_code", recorder.Status),
			attribute.Int64("http.response.body.size", recorder.Bytes),
			attribute.Int64("http.server.duration_ms", duration.Milliseconds()),
		)

		if recorder.Status >= http.StatusBadRequest {
			span.SetStatus(codes.Error, http.StatusText(recorder.Status))
			span.AddEvent(
				"http.request.failed",
				trace.WithAttributes(
					attribute.String("http.request.method", request.Method),
					attribute.String("http.route", route),
					attribute.String("url.path", request.URL.Path),
					attribute.Int("http.response.status_code", recorder.Status),
					attribute.Int64("http.response.body.size", recorder.Bytes),
				),
			)
		}

		coretelemetry.RecordHTTPServerRequest(
			ctx,
			edgeService.Key,
			request.Method,
			route,
			duration,
			recorder.Status,
		)

		coretelemetry.Logger(logger, ctx).Log(
			ctx,
			corehttp.LevelForStatus(recorder.Status),
			"request completed",
			slog.Group(
				"request",
				"method", request.Method,
				"route", route,
				"path", request.URL.Path,
				"contentLength", request.ContentLength,
			),
			slog.Group(
				"response",
				"status", recorder.Status,
				"bytes", recorder.Bytes,
			),
			"durationMs", duration.Milliseconds(),
		)
	})
}
