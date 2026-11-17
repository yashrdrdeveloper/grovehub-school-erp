package impl

import (
	config "grovehub/services/integration/internal/app/config"
	connectionpublic "grovehub/services/integration/internal/app/controller/connection/public"
	publicimpl "grovehub/services/integration/internal/app/controller/connection/public/impl"
	connectionprotected "grovehub/services/integration/internal/app/controller/connection/protected"
	protectedimpl "grovehub/services/integration/internal/app/controller/connection/protected/impl"
	providerpublic "grovehub/services/integration/internal/app/controller/provider/public"
	providerimpl "grovehub/services/integration/internal/app/controller/provider/public/impl"
	webhookrestricted "grovehub/services/integration/internal/app/controller/webhook/restricted"
	webhookimpl "grovehub/services/integration/internal/app/controller/webhook/restricted/impl"
	connections "grovehub/services/integration/internal/app/services/connection"
	providers "grovehub/services/integration/internal/app/services/provider"
)

// Each HTTP access tier owns a controller implementation; no cross-tier shared controller.
func (b *builder) newRoutes(connectionService connections.ConnectionService, providerService providers.ProviderService, push config.GmailPush) (connectionpublic.ConnectionController, connectionprotected.ConnectionController, providerpublic.ProviderController, webhookrestricted.WebhookController) {
	return publicimpl.New(connectionService),
		protectedimpl.New(connectionService),
		providerimpl.NewProviderController(providerService),
		webhookimpl.New(connectionService, push)
}
