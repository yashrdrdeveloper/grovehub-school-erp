package impl

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	coreerrors "grovehub/libs/core/errors"
	httprequest "grovehub/libs/core/http/request"
	controllerhelpers "grovehub/services/integration/internal/app/controller/helpers"
	connectionprotected "grovehub/services/integration/internal/app/controller/connection/protected"
	clientprincipal "grovehub/services/integration/client/principal"
	httpresponse "grovehub/libs/core/http/response"
	integrationerrors "grovehub/services/integration/internal/app/constants/error"
	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	providerdto "grovehub/services/integration/internal/app/providers/dto"
	integrationservices "grovehub/services/integration/internal/app/services/connection"
)

type Controller struct {
	service integrationservices.ConnectionService
}

var _ connectionprotected.ConnectionController = (*Controller)(nil)

func New(service integrationservices.ConnectionService) connectionprotected.ConnectionController {
	return &Controller{service: service}
}

func (c *Controller) principalFromRequest(request *http.Request) (clientprincipal.VerifiedPrincipal, bool) {
	return clientprincipal.From(request.Context())
}

func (c *Controller) writeServiceError(w http.ResponseWriter, r *http.Request, err error) {
	controllerhelpers.WriteServiceError(w, r, err)
}

func writeError(w http.ResponseWriter, r *http.Request, status int, definition coreerrors.Definition) {
	controllerhelpers.WriteError(w, r, status, definition)
}

func (c *Controller) owned(w http.ResponseWriter, r *http.Request) (clientprincipal.VerifiedPrincipal, bool) {
	p, ok := c.principalFromRequest(r)
	if !ok {
		writeError(w, r, http.StatusForbidden, integrationerrors.Forbidden)
		return p, false
	}
	current, err := c.service.Get(r.Context(), p.TenantID, r.PathValue("id"))
	if err != nil {
		c.writeServiceError(w, r, err)
		return p, false
	}
	if current.OwnerGID != p.ActorGID {
		writeError(w, r, http.StatusForbidden, integrationerrors.Forbidden)
		return p, false
	}
	return p, true
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	p, ok := c.principalFromRequest(r)
	if !ok {
		writeError(w, r, http.StatusForbidden, integrationerrors.Forbidden)
		return
	}
	var req connectiondto.CreateConnection
	if err := httprequest.DecodeStrictJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
		return
	}
	if req.OwnerGID != p.ActorGID {
		writeError(w, r, http.StatusForbidden, integrationerrors.Forbidden)
		return
	}
	result, err := c.service.Create(r.Context(), p.TenantID, req)
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusCreated, result)
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	p, ok := c.principalFromRequest(r)
	if !ok {
		writeError(w, r, http.StatusForbidden, integrationerrors.Forbidden)
		return
	}
	result, err := c.service.Get(r.Context(), p.TenantID, r.PathValue("id"))
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	if result.OwnerGID != p.ActorGID {
		writeError(w, r, http.StatusForbidden, integrationerrors.Forbidden)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusOK, result)
}

func (c *Controller) List(w http.ResponseWriter, r *http.Request) {
	p, ok := c.principalFromRequest(r)
	if !ok {
		writeError(w, r, http.StatusForbidden, integrationerrors.Forbidden)
		return
	}
	rows, err := c.service.List(r.Context(), p.TenantID, p.ActorGID)
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusOK, rows)
}

func (c *Controller) Disconnect(w http.ResponseWriter, r *http.Request) {
	p, ok := c.principalFromRequest(r)
	if !ok {
		writeError(w, r, http.StatusForbidden, integrationerrors.Forbidden)
		return
	}
	id := r.PathValue("id")
	current, err := c.service.Get(r.Context(), p.TenantID, id)
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	if current.OwnerGID != p.ActorGID {
		writeError(w, r, http.StatusForbidden, integrationerrors.Forbidden)
		return
	}
	result, err := c.service.Disconnect(r.Context(), p.TenantID, id)
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusOK, result)
}

func (c *Controller) BeginOAuth(w http.ResponseWriter, r *http.Request) {
	p, ok := c.owned(w, r)
	if !ok {
		return
	}
	result, err := c.service.BeginOAuth(r.Context(), p.TenantID, r.PathValue("id"))
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusOK, result)
}

func (c *Controller) RunTest(w http.ResponseWriter, r *http.Request) {
	p, ok := c.owned(w, r)
	if !ok {
		return
	}
	var payload struct {
		Direction string `json:"direction"`
	}
	if err := httprequest.DecodeStrictJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
		return
	}
	result, err := c.service.RunTest(r.Context(), p.TenantID, r.PathValue("id"), payload.Direction)
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusOK, result)
}

func (c *Controller) ListTests(w http.ResponseWriter, r *http.Request) {
	p, ok := c.owned(w, r)
	if !ok {
		return
	}
	result, err := c.service.ListTests(r.Context(), p.TenantID, r.PathValue("id"))
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusOK, result)
}

func (c *Controller) Watch(w http.ResponseWriter, r *http.Request) {
	p, ok := c.owned(w, r)
	if !ok {
		return
	}
	result, err := c.service.Watch(r.Context(), p.TenantID, r.PathValue("id"))
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusOK, result)
}

func (c *Controller) ExecuteAction(w http.ResponseWriter, r *http.Request) {
	p, ok := c.owned(w, r)
	if !ok {
		return
	}
	action := providerdto.Action(strings.TrimSpace(r.PathValue("action")))
	if action == "" || len(action) > 100 {
		writeError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
		return
	}
	var payload struct {
		Input json.RawMessage `json:"input"`
	}
	if err := httprequest.DecodeStrictJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
		return
	}
	result, err := c.service.ExecuteAction(r.Context(), p.TenantID, r.PathValue("id"), action, payload.Input)
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	_ = httpresponse.WriteSuccess(w, http.StatusOK, result)
}
