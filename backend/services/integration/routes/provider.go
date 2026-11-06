package routes

import (
 "net/http"
 "grovehub/services/integration/controller"
)
func registerProviderRoutes(mux *http.ServeMux,c controller.ProviderController){
 handle(mux,"GET","/integrations/providers",c.Providers)
 handle(mux,"GET","/integrations/providers/{provider}",c.Provider)
}
