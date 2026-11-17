package restricted

import "net/http"

// Webhook endpoints are workload-protected by the Launcher gateway.
// @Routes /webhooks
type WebhookController interface {
	// @POST /gmail
	GmailPush(http.ResponseWriter, *http.Request)
}
