package routes

import (
 "net/http"

 httptracing "grovehub/libs/core/http/tracing"
 "grovehub/services/integration/controller"
)
func RegisterHTTP(mux *http.ServeMux,c controller.ConnectionController) {
 routes:=[]struct{method,pattern string;handle http.HandlerFunc}{
  {"GET","/integrations/providers",c.Providers},
  {"POST","/integrations/connections",c.Create},
  {"GET","/integrations/connections",c.List},
  {"GET","/integrations/connections/{id}",c.Get},
  {"DELETE","/integrations/connections/{id}",c.Disconnect},
 }
 for _,route:=range routes {
  pattern:=route.method+" "+route.pattern
  mux.Handle(pattern,httptracing.ServiceHandler("integration",pattern,route.handle))
 }
}
