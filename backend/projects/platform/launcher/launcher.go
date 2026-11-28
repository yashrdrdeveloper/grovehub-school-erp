package launcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"grovehub/launcher/bootstrap"
	"grovehub/launcher/middleware"
	"grovehub/launcher/routing"
	"grovehub/launcher/services/gateway"
 schoolgateway "grovehub/launcher/services/gateway/session"
	coregateway "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	headercontext "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext"
	coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
	"grovehub/services/foundation"
	"grovehub/services/audit"
	"grovehub/services/identity"
	"grovehub/services/infrastructure"
	infrastructureconfig "grovehub/services/infrastructure/client/messaging"
	"grovehub/services/integration"
)

// Run assembles the backend modules behind one verified gateway.
func Run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var draining atomic.Bool
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

	tenantState,err:=gateway.NewTenantState(gateway.TenantAccessConfig{Origin:os.Getenv("PLATFORM_TENANT_ACCESS_ORIGIN"),KeyID:os.Getenv("PLATFORM_TENANT_ACCESS_KEY_ID"),PrivateKey:os.Getenv("PLATFORM_TENANT_ACCESS_PRIVATE_KEY")})
	if err!=nil{return err}
	tenantCheck:=schoolgateway.ActiveTenantCheck(tenantState)
 if os.Getenv("PLATFORM_TENANT_ACCESS_REQUIRED")=="true"&&tenantCheck==nil{return errors.New("tenant access verifier required")}
	// No loopback development principal is installed in the real launcher.
	integrationModule, err := integration.NewWithTenantState(
		logger,
		db,
		infrastructureModule.MessagingClient(),
		tenantState,
	)
	if err != nil {
		return fmt.Errorf("initialize Integration: %w", err)
	}
	defer integrationModule.Close()
	integrationModule.RegisterHTTP(mux)
	integration.RegisterMailReconciliationHTTP(mux, integration.NewMailReconciliationClient(db, identity.NewInvitationDeliveryReader(db)))

	identityModule, err := identity.New(
		logger,
		db,
		foundationModule.Client(),
		integrationModule.Client(),
		infrastructureModule.EventStore(),
		tenantState,
	)
	if err != nil {
		return fmt.Errorf("initialize Identity: %w", err)
	}
	defer identityModule.Close()
	identityModule.RegisterHTTP(mux)

	auditCtx,cancelAudit:=context.WithTimeout(context.Background(),5*time.Second)
	auditDB,closeAuditDB,err:=bootstrap.NewAuditPostgresDB(auditCtx,db)
	cancelAudit()
	if err!=nil { return err }
	defer closeAuditDB()
	auditModule, err := audit.New(logger, auditDB, infrastructureModule.MessagingClient(), identityModule)
	if err != nil { return fmt.Errorf("initialize Audit: %w", err) }
	defer auditModule.Close()
	auditModule.RegisterHTTP(mux)
	auditSQL,err:=auditDB.DB()
	if err!=nil{return err}
	routing.RegisterReadiness(mux,func() bool {
		if draining.Load() || !auditModule.Ready() { return false }
		probe,cancel:=context.WithTimeout(context.Background(),2*time.Second)
		defer cancel()
		if sqlDB.PingContext(probe)!=nil || auditSQL.PingContext(probe)!=nil { return false }
		return infrastructureModule.MessagingClient().Health(probe)==nil && !draining.Load()
	})

	routing.RegisterFallback(mux)

	// Both external HTTP access tiers and restricted provider webhooks
	// are verified before dispatch into a service-specific controller.
 if tenantState!=nil{
  stopReconciliation:=identity.StartTenantSessionReconciler(ctx,logger,db,tenantState)
  defer stopReconciliation()
 }
	workloads, err := gateway.NewPlatformWorkloads(os.Getenv("PLATFORM_WORKLOAD_CONFIG"), identityModule, integrationModule,tenantCheck)
	if err != nil { return fmt.Errorf("initialize workload gateway: %w", err) }
	gatewayRuntime := gateway.New(
		mux,
		identityModule,
		workloads,
		os.Getenv("IDENTITY_ALLOWED_ORIGIN"),
		nil, // Unconfigured privileged internal invocation fails closed.
		identityModule.RouteLimits(),
	)

	restrictedHandler := coregateway.RestrictedOnly(mux, gatewayRuntime.External)
	edgeService := coretelemetry.ServiceLauncherEdge
	edgeHandler := otelhttp.NewHandler(
		middleware.RequestObservability(
			logger,
			headercontext.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, pattern := mux.Handler(r)
				if pattern == "GET /health" || pattern == "GET /ready" {
					gatewayRuntime.External.ServeHTTP(w, r)
					return
				}
				restrictedHandler.ServeHTTP(w, r)
			})),
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
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 45 * time.Second,
		IdleTimeout: 60 * time.Second,
		MaxHeaderBytes: 1 << 20,
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

	done:=make(chan error,1)
	go func(){done<-server.Serve(listener)}()
	select {
	case err:=<-done:
		if err!=nil&&!errors.Is(err,http.ErrServerClosed){return fmt.Errorf("serve HTTP: %w",err)}
		return nil
	case <-ctx.Done():
		draining.Store(true)
		shutdown,cancel:=context.WithTimeout(context.Background(),15*time.Second)
		defer cancel()
		if err:=server.Shutdown(shutdown);err!=nil{_ = server.Close();return err}
		return nil
	}
}
