package helpers

import (
	"context"
	coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
	"net/http"
	"time"
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
