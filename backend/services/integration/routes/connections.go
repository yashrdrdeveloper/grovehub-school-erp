package routes

import (
 "net/http"

 "grovehub/services/integration/controller"
 routehelpers "grovehub/services/integration/routes/helpers"
)

// registerConnectionRoutes owns every /integrations/connections endpoint,
// including nested OAuth, tests, provider watches and action execution.
func registerConnectionRoutes(mux *http.ServeMux,c controller.Controller){
 routehelpers.Handle(mux,"POST","/integrations/connections",c.Create)
 routehelpers.Handle(mux,"GET","/integrations/connections",c.List)
 routehelpers.Handle(mux,"GET","/integrations/connections/{id}",c.Get)
 routehelpers.Handle(mux,"DELETE","/integrations/connections/{id}",c.Disconnect)

 routehelpers.Handle(mux,"POST","/integrations/connections/{id}/oauth/start",c.BeginOAuth)
 routehelpers.Handle(mux,"GET","/integrations/connections/oauth/{provider}/callback",c.FinishOAuth)

 routehelpers.Handle(mux,"POST","/integrations/connections/{id}/test",c.RunTest)
 routehelpers.Handle(mux,"GET","/integrations/connections/{id}/tests",c.ListTests)
 routehelpers.Handle(mux,"POST","/integrations/connections/{id}/watch",c.Watch)
 routehelpers.Handle(mux,"POST","/integrations/connections/{id}/actions/{action}",c.ExecuteAction)
}
