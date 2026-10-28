package impl

import (
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	httpresponse "grovehub/libs/core/http/response"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/controller"
	controllerhelpers "grovehub/services/mailProvider/controller/helpers"
	mailproviderrequest "grovehub/services/mailProvider/dto/request"
	mailproviderresponse "grovehub/services/mailProvider/dto/response"
	mailprovidermappers "grovehub/services/mailProvider/mappers"
	mailproviderservices "grovehub/services/mailProvider/services"
)

type SenderController struct {
	service mailproviderservices.SenderService
	mapper  mailprovidermappers.SenderMapper
}

func NewSenderController(
	service mailproviderservices.SenderService,
	mapper mailprovidermappers.SenderMapper,
) controller.SenderController {
	return &SenderController{service: service, mapper: mapper}
}

func (controller *SenderController) Register(writer http.ResponseWriter, request *http.Request) {
	var payload mailproviderrequest.RegisterSender
	if err := httprequest.DecodeJSON(request, &payload); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, mailprovidererror.InvalidSender)
		return
	}
	sender, err := controller.service.Register(request.Context(), mailproviderservices.RegisterSenderInput{
		DomainId: payload.DomainId, Email: payload.Email, DisplayName: payload.DisplayName,
	})
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusCreated, controller.mapper.ToResponse(sender))
}

func (controller *SenderController) ListAll(writer http.ResponseWriter, request *http.Request) {
	senders, err := controller.service.ListAll(request.Context())
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}

	response := make([]mailproviderresponse.Sender, 0, len(senders))
	for _, sender := range senders {
		response = append(response, controller.mapper.ToResponse(sender))
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, response)
}

func (controller *SenderController) GetById(writer http.ResponseWriter, request *http.Request) {
	sender, err := controller.service.GetById(request.Context(), request.PathValue("senderId"))
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, controller.mapper.ToResponse(sender))
}

func (controller *SenderController) Deactivate(writer http.ResponseWriter, request *http.Request) {
	sender, err := controller.service.Deactivate(request.Context(), request.PathValue("senderId"))
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, controller.mapper.ToResponse(sender))
}
