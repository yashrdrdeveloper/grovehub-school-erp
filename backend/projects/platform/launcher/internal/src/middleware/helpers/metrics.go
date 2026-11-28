package helpers

import (
	"context"
	"net/http"
	"time"

	coretelemetry "grovehub/libs/core/telemetry"
)

func RecordMetrics(ctx context.Context, request *http.Request, route string, recorder *Recorder, duration time.Duration) {
	edgeService := coretelemetry.ServiceLauncherEdge
	coretelemetry.RecordHTTPServerRequest(
		ctx,
		edgeService.Key,
		request.Method,
		route,
		duration,
		recorder.Status,
	)
}
