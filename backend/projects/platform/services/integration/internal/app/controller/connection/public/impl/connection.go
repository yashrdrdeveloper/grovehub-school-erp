package impl

import (
	"net/http"
	"strings"
	httpresponse "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/response"

		controllerhelpers "grovehub/services/integration/internal/app/controller/helpers"
	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	connectionpublic "grovehub/services/integration/internal/app/controller/connection/public"
	integrationerrors "grovehub/services/integration/internal/app/constants/error"
	integrationservices "grovehub/services/integration/internal/app/services/connection"
)

type Controller struct {
	service integrationservices.ConnectionService
}

var _ connectionpublic.ConnectionController = (*Controller)(nil)

func New(service integrationservices.ConnectionService) connectionpublic.ConnectionController {
	return &Controller{service: service}
}

func (c *Controller) writeServiceError(w http.ResponseWriter, r *http.Request, err error) {
	controllerhelpers.WriteServiceError(w, r, err)
}

func writeError(w http.ResponseWriter, r *http.Request, status int, definition coreerrors.Definition) {
	controllerhelpers.WriteError(w, r, status, definition)
}

func (c *Controller) FinishOAuth(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("error") != "" {
		writeError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
		return
	}
	result, err := c.service.FinishOAuth(r.Context(), strings.ToUpper(r.PathValue("provider")), r.URL.Query().Get("code"), r.URL.Query().Get("state"))
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusOK, result)
}
