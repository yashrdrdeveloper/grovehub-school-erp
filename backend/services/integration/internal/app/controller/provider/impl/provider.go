package impl

import (
    "net/http"

    httpresponse "grovehub/libs/core/http/response"
    controller "grovehub/services/integration/internal/app/controller"
    controllerhelpers "grovehub/services/integration/internal/app/controller/helpers"
    providerservice "grovehub/services/integration/internal/app/services/provider"
)

// ProviderController is independent of connections and owns public discovery.
type ProviderController struct { service providerservice.ProviderService }

var _ controller.ProviderController = (*ProviderController)(nil)

func NewProviderController(service providerservice.ProviderService) controller.ProviderController {
    return &ProviderController{service:service}
}

func (c *ProviderController) Providers(w http.ResponseWriter, r *http.Request) {
    _=httpresponse.WriteSuccess(w,http.StatusOK,c.service.Providers(r.Context()))
}

func (c *ProviderController) Provider(w http.ResponseWriter, r *http.Request) {
    result,err:=c.service.Provider(r.Context(),r.PathValue("provider"))
    if err!=nil {controllerhelpers.WriteServiceError(w,r,err);return}
    _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}
