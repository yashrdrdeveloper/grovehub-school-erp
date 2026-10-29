package routes

import (
	"net/http"

	httptracing "grovehub/libs/core/http/tracing"
	"grovehub/services/mailProvider/controller"
)

func RegisterHTTP(
	mux *http.ServeMux,
	domainController controller.DomainController,
	senderController controller.SenderController,
	messageController controller.MessageController,
) {
	register := func(pattern string, handler http.HandlerFunc) {
		mux.Handle(pattern, httptracing.ServiceHandler("mail-provider", pattern, handler))
	}

	register("POST /mail-provider/domains", domainController.Register)
	register("GET /mail-provider/domains", domainController.ListAll)
	register("GET /mail-provider/domains/by-name/{domainName}", domainController.GetByName)
	register("GET /mail-provider/domains/{domainId}", domainController.GetById)
	register("POST /mail-provider/domains/{domainId}/verify", domainController.Verify)
	register("POST /mail-provider/domains/{domainId}/deactivate", domainController.Deactivate)

	register("POST /mail-provider/senders", senderController.Register)
	register("GET /mail-provider/senders", senderController.ListAll)
	register("GET /mail-provider/senders/{senderId}", senderController.GetById)
	register("POST /mail-provider/senders/{senderId}/deactivate", senderController.Deactivate)

	register("POST /mail-provider/messages", messageController.Send)
	register("POST /mail-provider/messages/search", messageController.Search)
	register("POST /mail-provider/messages/{messageId}/retry", messageController.Retry)
	register("POST /mail-provider/messages/{messageId}/resend", messageController.Resend)
	register("POST /mail-provider/messages/{messageId}/reply", messageController.Reply)
	register("GET /mail-provider/messages/{messageId}", messageController.GetById)
}
