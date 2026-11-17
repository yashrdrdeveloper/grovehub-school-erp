package launcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"grovehub/launcher/bootstrap"
	"grovehub/launcher/middleware"
	"grovehub/launcher/routing"
	"grovehub/launcher/services/gateway"
	headercontext "grovehub/libs/core/headercontext"
	coretelemetry "grovehub/libs/core/telemetry"
	"grovehub/services/foundation"
	"grovehub/services/identity"
	"grovehub/services/infrastructure"
	infrastructureconfig "grovehub/services/infrastructure/client/messaging"
	"grovehub/services/integration"
)

// Run assembles the backend modules behind one verified gateway.
func Run(logger *slog.Logger) error {
	databaseCtx, cancelDatabase := context.WithTimeout(context.Background(), 5*time.Second)
	db, err := bootstrap.NewPostgresDB(databaseCtx)
	cancelDatabase()
	if err != nil {
		return fmt.Errorf("bootstrap PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access PostgreSQL connection: %w", err)
	}
	defer func() { _ = sqlDB.Close() }()

	infrastructureModule, err := infrastructure.New(infrastructureconfig.DefaultMessaging(), db, logger)
	if err != nil {
		return fmt.Errorf("initialize Infrastructure: %w", err)
	}
	defer func() { _ = infrastructureModule.Close() }()

	mux := http.NewServeMux()
	routing.RegisterHealth(mux)

	foundationModule, err := foundation.New(
		logger,
		db,
		infrastructureModule.EventStore(),
	)
	if err != nil {
		return fmt.Errorf("initialize Foundation: %w", err)
	}
	defer foundationModule.Close()
	foundationModule.RegisterHTTP(mux)

	// No loopback development principal is installed in the real launcher.
	integrationModule, err := integration.New(
		logger,
		db,
		infrastructureModule.MessagingClient(),
		integration.DevAuth{},
	)
	if err != nil {
		return fmt.Errorf("initialize Integration: %w", err)
	}
	defer integrationModule.Close()
	integrationModule.RegisterHTTP(mux)

	identityModule, err := identity.New(
		logger,
		db,
		foundationModule.Client(),
		integrationModule.Client(),
		infrastructureModule.EventStore(),
	)
	if err != nil {
		return fmt.Errorf("initialize Identity: %w", err)
	}
	defer identityModule.Close()
	identityModule.RegisterHTTP(mux)

	routing.RegisterFallback(mux)

	// Both external HTTP access tiers and restricted provider webhooks
	// are verified before dispatch into a service-specific controller.
	gatewayRuntime := gateway.New(
		mux,
		identityModule,
		integrationModule,
		os.Getenv("IDENTITY_ALLOWED_ORIGIN"),
		nil, // Unconfigured privileged internal invocation fails closed.
	)

	edgeService := coretelemetry.ServiceLauncherEdge
	edgeHandler := otelhttp.NewHandler(
		middleware.RequestObservability(
			logger,
			headercontext.Middleware(gatewayRuntime.External),
		),
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

	// Explicit deployment configuration is required for a non-loopback bind.
	httpAddress := os.Getenv("GROOVEHUB_HTTP_ADDR")
	if httpAddress == "" {
		httpAddress = "127.0.0.1:8080"
	}

	server := &http.Server{
		Addr: httpAddress,
		Handler: edgeHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	listener, err := net.Listen("tcp", httpAddress)
	if err != nil {
		return fmt.Errorf("bind HTTP listener %s: %w", httpAddress, err)
	}

	logger.Info(
		"HTTP server started",
		"service", edgeService.Name,
		"serviceKey", edgeService.Key,
		"address", listener.Addr().String(),
	)

	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	return nil
}
