package routes

import (
 "net/http"
 "grovehub/services/integration/controller"
)
func registerOAuthRoutes(mux *http.ServeMux,c controller.OAuthController){
 handle(mux,"POST","/integrations/connections/{id}/oauth/start",c.BeginOAuth)
 handle(mux,"GET","/integrations/oauth/{provider}/callback",c.FinishOAuth)
}
