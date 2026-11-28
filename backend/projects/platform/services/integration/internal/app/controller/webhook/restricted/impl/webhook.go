package impl

import (
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	controllerhelpers "grovehub/services/integration/internal/app/controller/helpers"
	integrationconfig "grovehub/services/integration/internal/app/config"
	integrationerrors "grovehub/services/integration/internal/app/constants/error"
	"grovehub/services/integration/internal/app/externalService/mail/gmail"
	webhookrestricted "grovehub/services/integration/internal/app/controller/webhook/restricted"
	integrationservices "grovehub/services/integration/internal/app/services/connection"
)

type Controller struct {
	service integrationservices.ConnectionService
	verifier *gmail.PushVerifier
	pushEnabled bool
}

var _ webhookrestricted.WebhookController = (*Controller)(nil)

func New(service integrationservices.ConnectionService, config integrationconfig.GmailPush) webhookrestricted.WebhookController {
	return &Controller{
		service: service,
		verifier: gmail.NewPushVerifier(config.Audience, config.ServiceAccount),
		pushEnabled: config.Enabled(),
	}
}

func writeError(w http.ResponseWriter, r *http.Request, status int, definition coreerrors.Definition) {
	controllerhelpers.WriteError(w, r, status, definition)
}

func (c *Controller) GmailPush(w http.ResponseWriter, r *http.Request) {
	if !c.pushEnabled {
		writeError(w, r, http.StatusServiceUnavailable, integrationerrors.NotReady)
		return
	}
	if err := c.verifier.Verify(r.Context(), r.Header.Get("Authorization")); err != nil {
		writeError(w, r, http.StatusUnauthorized, integrationerrors.Forbidden)
		return
	}
	data, err := gmail.DecodePush(r.Body)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, integrationerrors.Invalid)
		return
	}
	if err = c.service.ProcessPush(r.Context(), data); err != nil {
		// Pub/Sub retries on publish failure; Communication will ingest later.
		writeError(w, r, http.StatusServiceUnavailable, integrationerrors.Failed)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
