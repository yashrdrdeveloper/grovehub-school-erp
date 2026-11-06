package routes

import (
 "net/http"
 "grovehub/services/integration/controller"
)
func registerWebhookRoutes(mux *http.ServeMux,c controller.WebhookController){
 handle(mux,"POST","/integrations/webhooks/gmail",c.GmailPush)
}
