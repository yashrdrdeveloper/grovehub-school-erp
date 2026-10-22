package launcher

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"grovehub/services/identity"
)

const httpAddress = ":8080"

func Run(logger *slog.Logger) error {
	mux := http.NewServeMux()

	identityModule := identity.New(logger)
	identityModule.RegisterHTTP(mux)

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           requestLoggingMiddleware(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("http server started", "service", "platform", "address", httpAddress)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}

	return nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (recorder *statusRecorder) WriteHeader(status int) {
	recorder.status = status
	recorder.ResponseWriter.WriteHeader(status)
}

func requestLoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()
		recorder := &statusRecorder{ResponseWriter: writer, status: http.StatusOK}

		next.ServeHTTP(recorder, request)

		logger.Info(
			"request completed",
			"service", "platform",
			"method", request.Method,
			"path", request.URL.Path,
			"status", recorder.status,
			"durationMs", time.Since(startedAt).Milliseconds(),
		)
	})
}
