package launcher

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"grovehub/launcher/middleware"
	"grovehub/launcher/routing"
	"grovehub/services/identity"
)

const httpAddress = ":8080"

func Run(logger *slog.Logger) error {
	mux := http.NewServeMux()

	identityModule := identity.New(logger)
	identityModule.RegisterHTTP(mux)
	routing.RegisterFallback(mux)

	edgeHandler := otelhttp.NewHandler(
		middleware.RequestObservability(logger, mux),
		"grovehub.edge",
		otelhttp.WithSpanOptions(
			trace.WithAttributes(attribute.String("grovehub.component", "edge")),
		),
	)

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           edgeHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("http server started", "service", "platform", "address", httpAddress)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}

	return nil
}
