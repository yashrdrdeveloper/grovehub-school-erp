package routes

import (
 "net/http"
 "grovehub/services/integration/controller"
)
func registerWatchRoutes(mux *http.ServeMux,c controller.WatchController){
 handle(mux,"POST","/integrations/connections/{id}/watch",c.Watch)
}
