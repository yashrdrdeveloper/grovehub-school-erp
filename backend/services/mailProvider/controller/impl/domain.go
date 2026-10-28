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

type DomainController struct {
	service mailproviderservices.DomainService
	mapper  mailprovidermappers.DomainMapper
}

func NewDomainController(
	service mailproviderservices.DomainService,
	mapper mailprovidermappers.DomainMapper,
) controller.DomainController {
	return &DomainController{service: service, mapper: mapper}
}

func (controller *DomainController) Register(writer http.ResponseWriter, request *http.Request) {
	var payload mailproviderrequest.RegisterDomain
	if err := httprequest.DecodeJSON(request, &payload); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, mailprovidererror.InvalidDomain)
		return
	}
	domain, err := controller.service.Register(request.Context(), mailproviderservices.RegisterDomainInput{Name: payload.Name})
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusCreated, controller.mapper.ToResponse(domain))
}

func (controller *DomainController) ListAll(writer http.ResponseWriter, request *http.Request) {
	domains, err := controller.service.ListAll(request.Context())
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}

	response := make([]mailproviderresponse.Domain, 0, len(domains))
	for _, domain := range domains {
		response = append(response, controller.mapper.ToResponse(domain))
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, response)
}

func (controller *DomainController) GetById(writer http.ResponseWriter, request *http.Request) {
	domain, err := controller.service.GetById(request.Context(), request.PathValue("domainId"))
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, controller.mapper.ToResponse(domain))
}

func (controller *DomainController) GetByName(writer http.ResponseWriter, request *http.Request) {
	domain, err := controller.service.GetByName(request.Context(), request.PathValue("domainName"))
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, controller.mapper.ToResponse(domain))
}

func (controller *DomainController) Verify(writer http.ResponseWriter, request *http.Request) {
	domain, err := controller.service.Verify(request.Context(), request.PathValue("domainId"))
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, controller.mapper.ToResponse(domain))
}

func (controller *DomainController) Deactivate(writer http.ResponseWriter, request *http.Request) {
	domain, err := controller.service.Deactivate(request.Context(), request.PathValue("domainId"))
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, controller.mapper.ToResponse(domain))
}
