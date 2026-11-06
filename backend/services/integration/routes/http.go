package routes

import (
 "net/http"
 "grovehub/services/integration/controller"
)

// RegisterHTTP composes each capability route group once.
func RegisterHTTP(mux *http.ServeMux,c controller.Controller){
 registerProviderRoutes(mux,c)
 registerConnectionRoutes(mux,c)
 registerOAuthRoutes(mux,c)
 registerTestRoutes(mux,c)
 registerWatchRoutes(mux,c)
 registerWebhookRoutes(mux,c)
 registerActionRoutes(mux,c)
}
