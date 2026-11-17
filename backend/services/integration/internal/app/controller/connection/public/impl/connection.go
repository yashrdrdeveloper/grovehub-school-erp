package impl

import (
	"grovehub/services/integration/internal/app/controller/connection/public"
	shared "grovehub/services/integration/internal/app/controller/connection/impl"
)

type Controller struct {
	*shared.ConnectionController
}

var _ public.ConnectionController = (*Controller)(nil)

func New(sharedController *shared.ConnectionController) public.ConnectionController {
	return &Controller{ConnectionController: sharedController}
}
