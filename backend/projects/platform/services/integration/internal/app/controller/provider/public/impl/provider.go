package impl

import (
	"net/http"

	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
	providerpublic "grovehub/services/integration/internal/app/controller/provider/public"
	controllerhelpers "grovehub/services/integration/internal/app/controller/helpers"
	providerservice "grovehub/services/integration/internal/app/services/provider"
)

// ProviderController is independent of connections and owns public discovery.
type ProviderController struct {
	service providerservice.ProviderService
}

var _ providerpublic.ProviderController = (*ProviderController)(nil)

func NewProviderController(service providerservice.ProviderService) providerpublic.ProviderController {
	return &ProviderController{service: service}
}

func (controller *ProviderController) Providers(writer http.ResponseWriter, request *http.Request) {
	_ = corehttp.WriteSuccess(writer, http.StatusOK, controller.service.Providers(request.Context()))
}

func (controller *ProviderController) Provider(writer http.ResponseWriter, request *http.Request) {
	result, err := controller.service.Provider(request.Context(), request.PathValue("provider"))
	if err != nil {
		controllerhelpers.WriteServiceError(writer, request, err)
		return
	}
	_ = corehttp.WriteSuccess(writer, http.StatusOK, result)
}
