package impl

import (
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	httpresponse "grovehub/libs/core/http/response"
	infrastructureerror "grovehub/services/infrastructure/constants/error"
	"grovehub/services/infrastructure/controller"
	controllerhelpers "grovehub/services/infrastructure/controller/helpers"
	infrastructurerequest "grovehub/services/infrastructure/dto/request"
	infrastructureresponse "grovehub/services/infrastructure/dto/response"
	infrastructureservices "grovehub/services/infrastructure/services"
)

type DKIMController struct {
	service infrastructureservices.DKIMService
}

func NewDKIMController(service infrastructureservices.DKIMService) controller.DKIMController {
	return &DKIMController{service: service}
}

func (controller *DKIMController) Ensure(writer http.ResponseWriter, request *http.Request) {
	var payload infrastructurerequest.EnsureDKIM
	if err := httprequest.DecodeJSON(request, &payload); err != nil {
		_ = httpresponse.WriteErrorContext(
			request.Context(),
			writer,
			http.StatusBadRequest,
			infrastructureerror.InvalidDKIMIdentity,
		)
		return
	}

	identity, err := controller.service.Ensure(request.Context(), infrastructureservices.EnsureDKIMInput{
		Domain:   payload.Domain,
		Selector: payload.Selector,
	})
	if err != nil {
		controllerhelpers.WriteError(request.Context(), writer, err)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusOK, infrastructureresponse.DKIMIdentity{
		Domain:    identity.Domain,
		Selector:  identity.Selector,
		PublicKey: identity.PublicKey,
	})
}
