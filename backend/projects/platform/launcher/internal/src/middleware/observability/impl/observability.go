package impl

import (
 "grovehub/launcher/internal/src/middleware"
	"log/slog"
	"net/http"
	"time"
 "grovehub/launcher/internal/src/middleware/helpers"


	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"

	coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
)

type observability struct{logger *slog.Logger}
var _ middleware.Middleware=(*observability)(nil)
func New(logger *slog.Logger) middleware.Middleware {return &observability{logger:logger}}
func(o *observability) Wrap(next http.Handler) http.Handler {
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
        route:=helpers.Route(request,recorder.Status)
        helpers.RecordTrace(ctx,request,route,recorder,duration)
        helpers.RecordMetrics(ctx,request,route,recorder,duration)
        helpers.RecordLog(o.logger,ctx,request,route,recorder,duration)
    })
}
