package routes

import (
	"net/http"

	connectionpublic "grovehub/services/integration/internal/app/controller/connection/public"
	connectionprotected "grovehub/services/integration/internal/app/controller/connection/protected"
	providerpublic "grovehub/services/integration/internal/app/controller/provider/public"
	webhookrestricted "grovehub/services/integration/internal/app/controller/webhook/restricted"
)

// RegisterHTTP composes resource-owned annotated route groups.
// Gateway policy is enforced centrally in Launcher, not by the route generator.
// @Routes /integrations
func RegisterHTTP(mux *http.ServeMux, connectionsPublic connectionpublic.ConnectionController, connectionsProtected connectionprotected.ConnectionController, providers providerpublic.ProviderController, webhook webhookrestricted.WebhookController) {
	registerPublicConnectionRoutes(mux, connectionsPublic)
	registerProtectedConnectionRoutes(mux, connectionsProtected)
	registerPublicProviderRoutes(mux, providers)
	registerRestrictedWebhookRoutes(mux, webhook)
}
