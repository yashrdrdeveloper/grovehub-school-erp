package impl

import (
 "context"
 "errors"
 "fmt"
 "log/slog"
 "net/http"
 "os"
 "os/signal"
 "sync/atomic"
 "syscall"
 "time"
 "grovehub/launcher/internal/src/bootstrap"
 "grovehub/launcher/internal/src/bootstrap/helpers"
 "grovehub/launcher/internal/src/config"
 "grovehub/launcher/internal/src/routes"
 health "grovehub/launcher/internal/src/controller/impl"
 gateway "grovehub/launcher/internal/src/services/gateway/impl"
 "grovehub/launcher/internal/src/services/gateway/operations/workload"
 schoolgateway "grovehub/launcher/session"
 "grovehub/launcher/internal/src/server"
 serverimpl "grovehub/launcher/internal/src/server/impl"
 "grovehub/services/foundation"
 "grovehub/services/audit"
 "grovehub/services/identity"
 "grovehub/services/infrastructure"
 infrastructureconfig "grovehub/services/infrastructure/client/messaging"
 "grovehub/services/integration"
)

type builder struct {logger *slog.Logger; config config.Config}

var _ bootstrap.Bootstrap = (*builder)(nil)

func New(logger *slog.Logger) bootstrap.Bootstrap {
 if logger==nil {logger=slog.Default()}
 return &builder{logger:logger,config:config.FromEnv()}
}

func (b *builder) Build() (result bootstrap.Runtime, err error) {
 logger:=b.logger
 cleanup:=&helpers.Cleanup{}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	cleanup.Add(stop)
 defer func(){if err!=nil {cleanup.Close()}}()
	draining := &atomic.Bool{}
	databaseCtx, cancelDatabase := context.WithTimeout(context.Background(), 5*time.Second)
	db, err := NewPostgresDB(databaseCtx)
	cancelDatabase()
	if err != nil {
		return nil,fmt.Errorf("bootstrap PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil,fmt.Errorf("access PostgreSQL connection: %w", err)
	}
	cleanup.Add(func(){ _ = sqlDB.Close() })

	infrastructureModule, err := infrastructure.New(infrastructureconfig.DefaultMessaging(), db, logger)
	if err != nil {
		return nil,fmt.Errorf("initialize Infrastructure: %w", err)
	}
	cleanup.Add(func(){ _ = infrastructureModule.Close() })

	mux := http.NewServeMux()

	foundationModule, err := foundation.New(
		logger,
		db,
		infrastructureModule.EventStore(),
	)
	if err != nil {
		return nil,fmt.Errorf("initialize Foundation: %w", err)
	}
	cleanup.Add(foundationModule.Close)
	foundationModule.RegisterHTTP(mux)

	tenantState,err:=workload.NewTenantState(workload.TenantAccessConfig{Origin:b.config.TenantOrigin,KeyID:b.config.TenantKeyID,PrivateKey:b.config.TenantPrivateKey})
	if err!=nil{return nil,err}
	tenantCheck:=schoolgateway.ActiveTenantCheck(tenantState)
 if b.config.TenantRequired&&tenantCheck==nil{return nil,errors.New("tenant access verifier required")}
	// No loopback development principal is installed in the real launcher.
	integrationModule, err := integration.NewWithTenantState(
		logger,
		db,
		infrastructureModule.MessagingClient(),
		tenantState,
	)
	if err != nil {
		return nil,fmt.Errorf("initialize Integration: %w", err)
	}
	cleanup.Add(integrationModule.Close)
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
		return nil,fmt.Errorf("initialize Identity: %w", err)
	}
	cleanup.Add(identityModule.Close)
	identityModule.RegisterHTTP(mux)

	auditCtx,cancelAudit:=context.WithTimeout(context.Background(),5*time.Second)
	auditDB,closeAuditDB,err:=NewAuditPostgresDB(auditCtx,db)
	cancelAudit()
	if err!=nil { return nil,err }
	cleanup.Add(closeAuditDB)
	auditModule, err := audit.New(logger, auditDB, infrastructureModule.MessagingClient(), identityModule)
	if err != nil { return nil,fmt.Errorf("initialize Audit: %w", err) }
	cleanup.Add(auditModule.Close)
	auditModule.RegisterHTTP(mux)
	auditSQL,err:=auditDB.DB()
	if err!=nil{return nil,err}
	routes.RegisterHTTP(mux, health.New(func() bool {
		if draining.Load() || !auditModule.Ready() { return false }
		probe,cancel:=context.WithTimeout(context.Background(),2*time.Second)
		defer cancel()
		if sqlDB.PingContext(probe)!=nil || auditSQL.PingContext(probe)!=nil { return false }
		return infrastructureModule.MessagingClient().Health(probe)==nil && !draining.Load()
	}))

	routes.RegisterFallback(mux)

	// Both external HTTP access tiers and restricted provider webhooks
	// are verified before dispatch into a service-specific controller.
 if tenantState!=nil{
  stopReconciliation:=identity.StartTenantSessionReconciler(ctx,logger,db,tenantState)
  cleanup.Add(stopReconciliation)
 }
	workloads, err := workload.NewPlatformWorkloads(b.config.Workload, identityModule, integrationModule,tenantCheck)
	if err != nil { return nil,fmt.Errorf("initialize workload gateway: %w", err) }
	gatewayRuntime := gateway.New(
		mux,
		identityModule,
		workloads,
		b.config.AllowedOrigin,
		nil, // Unconfigured privileged internal invocation fails closed.
		identityModule.RouteLimits(),
	)


 edgeHandler:=newEdgeHandler(logger,mux,gatewayRuntime.External())
 var httpServer server.Server = serverimpl.New(b.config.Address,edgeHandler,logger,draining)
 return &runtime{ctx:ctx,server:httpServer,cleanup:cleanup},nil
}
