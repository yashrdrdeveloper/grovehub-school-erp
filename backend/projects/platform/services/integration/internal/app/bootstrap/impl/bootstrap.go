// Package impl assembles Integration without exposing its repositories/providers.
package impl

import (
    "context"
    "grovehub/libs/schoolgateway"
    "errors"
    "log/slog"
    "net/http"

    "gorm.io/gorm"

    coregateway "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
    gmail "grovehub/services/integration/internal/app/externalService/mail/gmail"

    integrationclient "grovehub/services/integration/client"
    integrationconfig "grovehub/services/integration/internal/app/config"
    bootstrap "grovehub/services/integration/internal/app/bootstrap"
    connectionpublic "grovehub/services/integration/internal/app/controller/connection/public"
    connectionprotected "grovehub/services/integration/internal/app/controller/connection/protected"
    providerpublic "grovehub/services/integration/internal/app/controller/provider/public"
    webhookrestricted "grovehub/services/integration/internal/app/controller/webhook/restricted"
    "grovehub/services/integration/internal/app/routes"
    messaging "grovehub/services/infrastructure/client/messaging"
)

type builder struct {
    logger *slog.Logger
    db *gorm.DB
    messaging messaging.MessagingClient
    devAuth integrationconfig.DevAuth
    tenantState schoolgateway.TenantStateLookup
}

var _ bootstrap.Bootstrap = (*builder)(nil)

func New(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, auth integrationconfig.DevAuth, lookups ...schoolgateway.TenantStateLookup) bootstrap.Bootstrap {
    b:=&builder{logger:logger,db:db,messaging:events,devAuth:auth}
    if len(lookups)==1 {b.tenantState=lookups[0]}
    return b
}

type runtime struct {
    logger *slog.Logger
    connectionPublic connectionpublic.ConnectionController
    connectionProtected connectionprotected.ConnectionController
    providerPublic providerpublic.ProviderController
    webhookRestricted webhookrestricted.WebhookController
    client integrationclient.Client
    localOnly bool
    webhookVerifier *gmail.PushVerifier
    webhookEnabled bool
    stop context.CancelFunc
    done chan struct{}
}

var _ bootstrap.Runtime = (*runtime)(nil)

func (b *builder) Build() (bootstrap.Runtime,error) {
    if b.db==nil {return nil,errors.New("integration bootstrap requires a database")}
    if b.messaging==nil {return nil,errors.New("integration bootstrap requires Infrastructure messaging")}
    if b.logger==nil {b.logger=slog.Default()}
    if b.devAuth.Enabled {b.logger.Warn("Integration development authentication bypass enabled for loopback clients only; never enable in production")}
    providers,err:=b.newProviders()
    if err!=nil{return nil,err}
    eventProducer:=b.newEvents()
    connections,providerService,err:=b.newServices(providers,eventProducer)
    if err!=nil{return nil,err}
    connectionsPublic, connectionsProtected, providersPublic, webhooksRestricted := b.newRoutes(connections, providerService, providers.push)
    client:=b.newClient(connections,providerService)
    stop,done:=b.newWatchRenewal(connections,providers.push.Enabled()&&providers.gmail!=nil)
    return &runtime{logger:b.logger,connectionPublic:connectionsPublic, connectionProtected:connectionsProtected, providerPublic:providersPublic, webhookRestricted:webhooksRestricted,client:client,localOnly:b.devAuth.Enabled,webhookVerifier:gmail.NewPushVerifier(providers.push.Audience,providers.push.ServiceAccount),webhookEnabled:providers.push.Enabled(),stop:stop,done:done},nil
}

func (module *runtime) RegisterHTTP(mux *http.ServeMux) {routes.RegisterHTTP(mux,module.connectionPublic,module.connectionProtected,module.providerPublic,module.webhookRestricted)}

func (module *runtime) Client() integrationclient.Client {return module.client}

func (module *runtime) LocalOnly() bool {return module.localOnly}

func (module *runtime) Close() {module.stop();<-module.done}

var _ coregateway.WorkloadVerifier = (*runtime)(nil)

// VerifyWorkload authenticates the configured Gmail Pub/Sub push caller.
// Other restricted paths fail closed until an explicit workload strategy exists.
func (module *runtime) VerifyWorkload(ctx context.Context, request *http.Request) (context.Context, error) {
	if !module.webhookEnabled || module.webhookVerifier == nil ||
		request.Method != http.MethodPost ||
		request.URL.Path != "/integrations/restricted/webhooks/gmail" {
		return nil, errors.New("unrecognized restricted workload")
	}
	if err := module.webhookVerifier.Verify(ctx, request.Header.Get("Authorization")); err != nil {
		return nil, err
	}
	return ctx, nil
}
