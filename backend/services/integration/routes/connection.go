package routes

import (
 "net/http"
 "grovehub/services/integration/controller"
)
func registerConnectionRoutes(mux *http.ServeMux,c controller.ConnectionController){
 handle(mux,"POST","/integrations/connections",c.Create)
 handle(mux,"GET","/integrations/connections",c.List)
 handle(mux,"GET","/integrations/connections/{id}",c.Get)
 handle(mux,"DELETE","/integrations/connections/{id}",c.Disconnect)
}
