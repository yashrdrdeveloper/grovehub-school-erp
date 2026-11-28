package helpers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	corehttp "grovehub/libs/core/http"
	coretelemetry "grovehub/libs/core/telemetry"
)

func RecordLog(logger *slog.Logger, ctx context.Context, request *http.Request, route string, recorder *Recorder, duration time.Duration) {
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
}
