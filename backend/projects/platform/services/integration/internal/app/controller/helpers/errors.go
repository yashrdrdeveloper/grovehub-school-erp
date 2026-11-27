package helpers

import (
    providerhelpers "grovehub/services/integration/internal/app/providers/helpers"
    "errors"
    "net/http"

    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    httpresponse "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/response"
    integrationerrors "grovehub/services/integration/internal/app/constants/error"
    integrationrepository "grovehub/services/integration/internal/app/repository/connection"
    serviceimpl "grovehub/services/integration/internal/app/services/connection/impl"
)

// WriteServiceError preserves Integration's existing structured HTTP error mapping
// for both Connection and Provider controllers.
func WriteServiceError(w http.ResponseWriter, r *http.Request, err error) {
	if status := providerhelpers.HTTPStatus(err); status != 0 {
		switch status {
		case http.StatusBadRequest:
			WriteError(w, r, status, integrationerrors.Invalid)
		case http.StatusUnauthorized, http.StatusForbidden:
			WriteError(w, r, status, integrationerrors.Forbidden)
		case http.StatusNotFound:
			WriteError(w, r, status, integrationerrors.NotFound)
		case http.StatusConflict:
			WriteError(w, r, status, integrationerrors.NotReady)
		default:
			WriteError(w, r, status, integrationerrors.Failed)
		}
		return
	}
	switch {
	case errors.Is(err, serviceimpl.ErrInvalidConnection), errors.Is(err, integrationrepository.ErrConnectionNotFound):
		if errors.Is(err, integrationrepository.ErrConnectionNotFound) {
			WriteError(w, r, http.StatusNotFound, integrationerrors.NotFound)
		} else {
			WriteError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
		}
	case errors.Is(err, providerhelpers.ErrInvalidActionInput):
		WriteError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
	case errors.Is(err, providerhelpers.ErrUnsupportedAction):
		WriteError(w, r, http.StatusNotImplemented, integrationerrors.NotReady)
	case errors.Is(err, providerhelpers.ErrProviderUnavailable):
		WriteError(w, r, http.StatusServiceUnavailable, integrationerrors.NotReady)
	case errors.Is(err, providerhelpers.ErrUnknownProvider):
		WriteError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
	case errors.Is(err, serviceimpl.ErrOAuthInvalid):
		WriteError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
	case errors.Is(err, serviceimpl.ErrPushUnavailable):
		WriteError(w, r, http.StatusServiceUnavailable, integrationerrors.NotReady)
	case errors.Is(err, serviceimpl.ErrConnectionInactive):
		WriteError(w, r, http.StatusConflict, integrationerrors.NotReady)
	case errors.Is(err, serviceimpl.ErrOAuthUnavailable), errors.Is(err, serviceimpl.ErrConnectionNotReady):
		WriteError(w, r, http.StatusNotImplemented, integrationerrors.NotReady)
	default:
		WriteError(w, r, http.StatusInternalServerError, integrationerrors.Failed)
	}
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, definition coreerrors.Definition) {
	_ = httpresponse.WriteErrorContext(r.Context(), w, status, definition)
}
