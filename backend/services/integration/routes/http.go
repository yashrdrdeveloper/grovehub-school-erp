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
  {"POST","/integrations/connections/{id}/oauth/start",c.BeginOAuth},
  {"GET","/integrations/oauth/gmail/callback",c.FinishOAuth},
  {"POST","/integrations/connections/{id}/test",c.RunTest},
  {"GET","/integrations/connections/{id}/tests",c.ListTests},
  {"POST","/integrations/connections/{id}/sync",c.Sync},
  {"POST","/integrations/connections/{id}/watch",c.Watch},
  {"POST","/integrations/webhooks/gmail",c.GmailPush},
 }
 for _,route:=range routes {
  pattern:=route.method+" "+route.pattern
  mux.Handle(pattern,httptracing.ServiceHandler("integration",pattern,route.handle))
 }
}
