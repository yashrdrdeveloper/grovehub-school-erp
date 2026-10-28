package impl

import (
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	httpresponse "grovehub/libs/core/http/response"
	coresearch "grovehub/libs/core/search"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/controller"
	controllerhelpers "grovehub/services/mailProvider/controller/helpers"
	mailproviderrequest "grovehub/services/mailProvider/dto/request"
	mailproviderresponse "grovehub/services/mailProvider/dto/response"
	mailprovidermappers "grovehub/services/mailProvider/mappers"
	mailproviderservices "grovehub/services/mailProvider/services"
)

type MessageController struct {
	service mailproviderservices.MessageService
	mapper  mailprovidermappers.MessageMapper
}

func NewMessageController(
	service mailproviderservices.MessageService,
	mapper mailprovidermappers.MessageMapper,
) controller.MessageController {
	return &MessageController{service: service, mapper: mapper}
}

func (controller *MessageController) Send(writer http.ResponseWriter, request *http.Request) {
	var payload mailproviderrequest.SendMessage
	if err := httprequest.DecodeJSON(request, &payload); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, mailprovidererror.InvalidMessage)
		return
	}
	message, err := controller.service.Send(request.Context(), mailproviderservices.SendMessageInput{
		GID: payload.GID, SenderId: payload.SenderId, To: payload.To, Cc: payload.Cc, Bcc: payload.Bcc,
		Subject: payload.Subject, TextBody: payload.TextBody, HTMLBody: payload.HTMLBody,
	})
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusCreated, controller.mapper.ToResponse(message))
}

func (controller *MessageController) Search(writer http.ResponseWriter, request *http.Request) {
	var searchRequest coresearch.Request
	if err := httprequest.DecodeJSON(request, &searchRequest); err != nil {
		_ = httpresponse.WriteErrorContext(
			request.Context(),
			writer,
			http.StatusBadRequest,
			mailprovidererror.InvalidMessageSearch,
		)
		return
	}

	result, err := controller.service.Search(request.Context(), searchRequest)
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}

	items := make([]mailproviderresponse.Message, 0, len(result.Items))
	for _, message := range result.Items {
		items = append(items, controller.mapper.ToResponse(message))
	}
	response := coresearch.Result[mailproviderresponse.Message]{
		Items:    items,
		PageInfo: result.PageInfo,
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, response)
}

func (controller *MessageController) GetById(writer http.ResponseWriter, request *http.Request) {
	message, err := controller.service.GetById(request.Context(), request.PathValue("messageId"))
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, controller.mapper.ToResponse(message))
}
