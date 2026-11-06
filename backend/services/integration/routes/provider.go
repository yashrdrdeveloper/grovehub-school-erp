package routes

import (
	"net/http"

	"grovehub/services/integration/controller"
	routehelpers "grovehub/services/integration/routes/helpers"
)

func registerProviderRoutes(mux *http.ServeMux, c controller.ProviderController) {
	routehelpers.Handle(mux, "GET", "/integrations/providers", c.Providers)
	routehelpers.Handle(mux, "GET", "/integrations/providers/{provider}", c.Provider)
}
