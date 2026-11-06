package integration

import (
	"context"
	integrationservices "grovehub/services/integration/services"
	"log/slog"
	"net/http"
	"time"

	"gorm.io/gorm"

	infrastructureclient "grovehub/services/infrastructure/client"
	integrationclient "grovehub/services/integration/client"
	clientimpl "grovehub/services/integration/client/impl"
	configimpl "grovehub/services/integration/config/impl"
	"grovehub/services/integration/controller"
	controllerimpl "grovehub/services/integration/controller/impl"
	mailfactory "grovehub/services/integration/providers/mailprovider/impl"
	"grovehub/services/integration/providers/mailprovider/providers/gmail"
	"grovehub/services/integration/providers/mailprovider/providers/outlook"
	repositoryimpl "grovehub/services/integration/repository/impl"
	"grovehub/services/integration/routes"
	serviceimpl "grovehub/services/integration/services/impl"
)

type Module interface {
	RegisterHTTP(*http.ServeMux)
	Client() integrationclient.Client
	Close()
}
type module struct {
	controller controller.Controller
	client     integrationclient.Client
	stop       chan struct{}
}

func New(logger *slog.Logger, db *gorm.DB, messaging infrastructureclient.MessagingClient) (Module, error) {
	devAuth,devErr:=configimpl.NewEnvironmentLoader().DevAuth()
 if devErr!=nil{return nil,devErr}
 if devAuth.Enabled{
  logger.Warn("Integration development authentication bypass enabled for loopback clients only; never enable in production")
 }
 repository := repositoryimpl.NewConnectionRepository(db)
	loader := configimpl.NewEnvironmentLoader()
	gmailConfig, configErr := loader.Gmail()
	var client *gmail.Client
	key, keyErr := loader.CredentialKey()
	if keyErr != nil {
		logger.Warn("integration encryption key not configured; OAuth providers disabled")
	}
	if configErr == nil && keyErr == nil {
		client = gmail.NewClient(gmailConfig)
	} else {
		logger.Warn("gmail provider not configured; OAuth and tests disabled")
	}
	outlookConfig, outlookErr := loader.Outlook()
	var outlookClient *outlook.Client
	if outlookErr == nil && keyErr == nil {
		outlookClient = outlook.NewClient(outlookConfig)
	} else {
		logger.Warn("Outlook OAuth not configured; Outlook actions disabled")
	}
	runtime, err := mailfactory.NewFactory().Build(
		gmail.NewStrategy(client),
		outlook.NewStrategy(outlookClient),
	)
	if err != nil {
		return nil, err
	}
	registry, actions, oauth := runtime.Registry, runtime.Actions, runtime.OAuth
	service, err := serviceimpl.NewConnectionService(repository, repositoryimpl.NewWatchRepository(db), registry, actions, oauth, client, outlookClient, key, loader.GmailPush().Topic, messaging)
	if err != nil {
		return nil, err
	}
	push := loader.GmailPush()
	providerService := serviceimpl.NewProviderService(registry, actions, oauth)
	mod := &module{
		controller: controllerimpl.NewConnectionController(service, providerService, push, devAuth),
		client:     clientimpl.NewClient(service, providerService),
		stop:       make(chan struct{}),
	}
	if push.Enabled() && client != nil {
		go mod.watchRenewal(service)
	}
	return mod, nil
}
func (m *module) RegisterHTTP(mux *http.ServeMux)  { routes.RegisterHTTP(mux, m.controller) }
func (m *module) Client() integrationclient.Client { return m.client }

func (m *module) Close() { close(m.stop) }
func (m *module) watchRenewal(service integrationservices.ConnectionService) {
	// Distributed duplicate renewal is safe: Gmail watch registration is
	// idempotent and the DB upsert keeps only the latest expiration.
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-m.stop:
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			_ = service.RenewExpiring(ctx)
			cancel()
		}
	}
}
