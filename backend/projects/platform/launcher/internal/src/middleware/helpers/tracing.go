package helpers

import (
 "context";"net/http";"time"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/codes"
 "go.opentelemetry.io/otel/trace"
 coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
)
func RecordTrace(ctx context.Context,request *http.Request,route string,recorder *Recorder,duration time.Duration){
 edgeService:=coretelemetry.ServiceLauncherEdge
 span:=trace.SpanFromContext(ctx)
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


}
