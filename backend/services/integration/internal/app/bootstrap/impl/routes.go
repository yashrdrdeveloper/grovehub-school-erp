package impl

import (
	integrationconfig "grovehub/services/integration/internal/app/config"
	connectionpublic "grovehub/services/integration/internal/app/controller/connection/public"
	connectionpublicimpl "grovehub/services/integration/internal/app/controller/connection/public/impl"
	connectionprotected "grovehub/services/integration/internal/app/controller/connection/protected"
	connectionprotectedimpl "grovehub/services/integration/internal/app/controller/connection/protected/impl"
	providerpublic "grovehub/services/integration/internal/app/controller/provider/public"
	providerpublicimpl "grovehub/services/integration/internal/app/controller/provider/public/impl"
	webhookrestricted "grovehub/services/integration/internal/app/controller/webhook/restricted"
	webhookrestrictedimpl "grovehub/services/integration/internal/app/controller/webhook/restricted/impl"
	connectionservice "grovehub/services/integration/internal/app/services/connection"
	providerservice "grovehub/services/integration/internal/app/services/provider"
)

// Each access tier owns its controller implementation. No cross-tier shared
// controller implementation or redundant forwarding wrapper is required.
func (b *builder) newRoutes(connections connectionservice.ConnectionService, providers providerservice.ProviderService, push integrationconfig.GmailPush) (connectionpublic.ConnectionController, connectionprotected.ConnectionController, providerpublic.ProviderController, webhookrestricted.WebhookController) {
	return connectionpublicimpl.New(connections),
		connectionprotectedimpl.New(connections),
		providerpublicimpl.NewProviderController(providers),
		webhookrestrictedimpl.New(connections, push)
}
