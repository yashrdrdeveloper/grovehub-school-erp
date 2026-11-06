package routes

import (
 "net/http"

 "grovehub/services/integration/controller"
 routehelpers "grovehub/services/integration/routes/helpers"
)

func registerWebhookRoutes(mux *http.ServeMux,c controller.WebhookController){
 routehelpers.Handle(mux,"POST","/integrations/webhooks/gmail",c.GmailPush)
}
