package launcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
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
	headercontext "grovehub/libs/core/headercontext"
	"grovehub/services/foundation"
	"grovehub/services/identity"
 "grovehub/services/integration"
 integrationconfig "grovehub/services/integration/config"
	"grovehub/services/infrastructure"
	infrastructureconfig "grovehub/services/infrastructure/client/messaging"
)

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

	infrastructureModule, err := infrastructure.New(infrastructureconfig.DefaultMessaging(), db, logger)
	if err != nil {
		return fmt.Errorf("bootstrap infrastructure messaging: %w", err)
	}
	defer func() { _ = infrastructureModule.Close() }()

	mux := http.NewServeMux()
	routing.RegisterHealth(mux)

	foundationModule, err := foundation.New(
        logger,
        db,
        infrastructureModule.EventStore(),
    )
    if err != nil { return fmt.Errorf("bootstrap foundation: %w", err) }
    defer foundationModule.Close()
    foundationModule.RegisterHTTP(mux)

	identityModule := identity.New(
		logger,
		db,
		foundationModule.Client(),
		infrastructureModule.MessagingClient(),
	)
	identityModule.RegisterHTTP(mux)

 // TEMPORARY: local-only Integration identity until IAM middleware exists.
 // All protected Integration routes use this fixed development principal.
 // Remove when Identity provides a verified principal; never deploy this mode.
 integrationDevAuth := integrationconfig.DevAuth{
  Enabled: true,
  TenantID: "018f5d04-315c-75a9-b890-df0d02e787c1",
  ActorGID: "01a1b210-402c-7622-a33c-f82ceba26c7d",
 }
 integrationModule, err := integration.New(logger,db,infrastructureModule.MessagingClient(),integrationDevAuth)
 if err != nil { return fmt.Errorf("initialize integration module: %w",err) }
 integrationModule.RegisterHTTP(mux)
 defer integrationModule.Close()

 // A development principal must never expose the HTTP listener to the LAN.
 // Bind to loopback even if the default platform listener binds all interfaces.
 httpAddress := ":8080"
 if integrationModule.LocalOnly(){
  httpAddress="127.0.0.1:8080"
 }


	routing.RegisterFallback(mux)

	edgeService := coretelemetry.ServiceLauncherEdge
	edgeHandler := otelhttp.NewHandler(
		middleware.RequestObservability(logger, headercontext.Middleware(mux)),
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

	listener, err := net.Listen("tcp", httpAddress)
	if err != nil {
		return fmt.Errorf("bind HTTP listener %s: %w", httpAddress, err)
	}

	logger.Info(
		"http server started",
		"service", edgeService.Name,
		"serviceKey", edgeService.Key,
		"address", listener.Addr().String(),
	)

	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}

	return nil
}
