package routes

import (
 "net/http"
 "grovehub/services/integration/controller"
)
func registerTestRoutes(mux *http.ServeMux,c controller.TestController){
 handle(mux,"POST","/integrations/connections/{id}/test",c.RunTest)
 handle(mux,"GET","/integrations/connections/{id}/tests",c.ListTests)
}
