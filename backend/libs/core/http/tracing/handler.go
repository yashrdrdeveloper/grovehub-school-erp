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
		targetService := coretelemetry.ResolveService(service)
		callerService, ok := coretelemetry.ServiceFromContext(request.Context())
		if !ok {
			callerService = coretelemetry.ServiceLauncherEdge
		}

		clientContext, clientSpan := coretelemetry.StartSpan(
			request.Context(),
			callerService.Key,
			"call "+targetService.Key+" "+operation,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(
				attribute.String("grovehub.boundary", "in-process-service"),
				attribute.String("grovehub.target.service.key", targetService.Key),
				attribute.String("grovehub.target.service.name", targetService.Name),
				attribute.String("peer.service", targetService.Name),
			),
		)
		defer clientSpan.End()

		ctx, span := coretelemetry.StartSpan(
			clientContext,
			targetService.Key,
			operation,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("grovehub.boundary", "in-process-service"),
				attribute.String("grovehub.caller.service.key", callerService.Key),
				attribute.String("grovehub.caller.service.name", callerService.Name),
			),
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
		clientSpan.SetAttributes(
			attribute.String("http.request.method", request.Method),
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
		clientSpan.SetAttributes(
			attribute.Int("http.response.status_code", recorder.status),
		)
		if recorder.status >= http.StatusBadRequest {
			span.SetStatus(codes.Error, http.StatusText(recorder.status))
			clientSpan.SetStatus(codes.Error, http.StatusText(recorder.status))
			span.AddEvent(
				"service.operation.failed",
				trace.WithAttributes(
					attribute.String("grovehub.operation", operation),
					attribute.Int("http.response.status_code", recorder.status),
				),
			)
			clientSpan.AddEvent(
				"service.call.failed",
				trace.WithAttributes(
					attribute.String("grovehub.target.service.key", targetService.Key),
					attribute.Int("http.response.status_code", recorder.status),
				),
			)
		}

		coretelemetry.RecordOperation(
			ctx,
			targetService.Key,
			operation,
			time.Since(startedAt),
			recorder.status,
		)
	})
}
