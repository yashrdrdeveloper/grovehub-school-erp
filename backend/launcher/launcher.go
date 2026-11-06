package launcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"grovehub/launcher/bootstrap"
	"grovehub/launcher/middleware"
	"grovehub/launcher/routing"
	coretelemetry "grovehub/libs/core/telemetry"
	"grovehub/services/foundation"
	"grovehub/services/identity"
 "grovehub/services/integration"
	"grovehub/services/infrastructure"
	infrastructureconfig "grovehub/services/infrastructure/config"
)

const httpAddress = ":8080"

func Run(logger *slog.Logger) error {
	databaseContext, cancelDatabase := context.WithTimeout(context.Background(), 5*time.Second)
	db, err := bootstrap.NewPostgresDB(databaseContext)
	cancelDatabase()
	if err != nil {
		return fmt.Errorf("bootstrap postgres gorm db: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access postgres gorm db: %w", err)
	}
	defer func() { _ = sqlDB.Close() }()

	infrastructureModule, err := infrastructure.New(infrastructureconfig.DefaultMessaging())
	if err != nil {
		return fmt.Errorf("bootstrap infrastructure messaging: %w", err)
	}
	defer func() { _ = infrastructureModule.Close() }()

	mux := http.NewServeMux()
	routing.RegisterHealth(mux)

	foundationModule := foundation.New(
		logger,
		db,
		infrastructureModule.MessagingClient(),
	)
	foundationModule.RegisterHTTP(mux)

	identityModule := identity.New(
		logger,
		db,
		foundationModule.Client(),
		infrastructureModule.MessagingClient(),
	)
	identityModule.RegisterHTTP(mux)

 integrationModule, err := integration.New(logger,db)
 if err != nil { return fmt.Errorf("initialize integration module: %w",err) }
 integrationModule.RegisterHTTP(mux)
 defer integrationModule.Close()

	routing.RegisterFallback(mux)

	edgeService := coretelemetry.ServiceLauncherEdge
	edgeHandler := otelhttp.NewHandler(
		middleware.RequestObservability(logger, mux),
		"grovehub.edge",
		otelhttp.WithTracerProvider(coretelemetry.TracerProvider(edgeService.Key)),
		otelhttp.WithMeterProvider(coretelemetry.MeterProvider(edgeService.Key)),
		otelhttp.WithSpanNameFormatter(func(_ string, request *http.Request) string {
			route := strings.TrimPrefix(request.Pattern, request.Method+" ")
			if route == "" {
				route = "unmatched"
			}
			return edgeService.Name + " " + request.Method + " " + route
		}),
		otelhttp.WithSpanOptions(
			trace.WithAttributes(
				attribute.String("grovehub.service.key", edgeService.Key),
				attribute.String("grovehub.service.display_name", edgeService.Name),
			),
		),
	)

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           edgeHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info(
		"http server started",
		"service", edgeService.Name,
		"serviceKey", edgeService.Key,
		"address", httpAddress,
	)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}

	return nil
}
