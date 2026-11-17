package impl

import (
	"grovehub/services/integration/internal/app/controller/webhook/restricted"
	shared "grovehub/services/integration/internal/app/controller/connection/impl"
)

type Controller struct {
	*shared.ConnectionController
}

var _ restricted.WebhookController = (*Controller)(nil)

func New(sharedController *shared.ConnectionController) restricted.WebhookController {
	return &Controller{ConnectionController: sharedController}
}
