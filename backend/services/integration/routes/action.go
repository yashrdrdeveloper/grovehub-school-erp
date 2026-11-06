package routes

import (
 "net/http"
 "grovehub/services/integration/controller"
)
func registerActionRoutes(mux *http.ServeMux,c controller.ActionController){
 handle(mux,"POST","/integrations/connections/{id}/actions/{action}",c.ExecuteAction)
}
