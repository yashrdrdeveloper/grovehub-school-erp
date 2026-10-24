package impl

import (
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	httpresponse "grovehub/libs/core/http/response"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/controller"
	controllerhelpers "grovehub/services/mailProvider/controller/helpers"
	mailproviderrequest "grovehub/services/mailProvider/dto/request"
	mailprovidermappers "grovehub/services/mailProvider/mappers"
	mailproviderservices "grovehub/services/mailProvider/services"
)

type MessageController struct {
	service mailproviderservices.MessageService
	mapper  mailprovidermappers.MessageMapper
}

func NewMessageController(service mailproviderservices.MessageService, mapper mailprovidermappers.MessageMapper) controller.MessageController {
	return &MessageController{service: service, mapper: mapper}
}

func (controller *MessageController) Send(writer http.ResponseWriter, request *http.Request) {
	var payload mailproviderrequest.SendMessage
	if err := httprequest.DecodeJSON(request, &payload); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, mailprovidererror.InvalidMessage)
		return
	}
	message, err := controller.service.Send(request.Context(), mailproviderservices.SendMessageInput{
		SenderId: payload.SenderId, To: payload.To, Cc: payload.Cc, Bcc: payload.Bcc,
		Subject: payload.Subject, TextBody: payload.TextBody, HTMLBody: payload.HTMLBody,
	})
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusCreated, controller.mapper.ToResponse(message))
}

func (controller *MessageController) GetById(writer http.ResponseWriter, request *http.Request) {
	message, err := controller.service.GetById(request.Context(), request.PathValue("messageId"))
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, controller.mapper.ToResponse(message))
}
