package impl

import (
	"grovehub/services/integration/internal/app/controller/connection/protected"
	shared "grovehub/services/integration/internal/app/controller/connection/impl"
)

// Controller presents only authenticated connection operations.
type Controller struct {
	*shared.ConnectionController
}

var _ protected.ConnectionController = (*Controller)(nil)

func New(sharedController *shared.ConnectionController) protected.ConnectionController {
	return &Controller{ConnectionController: sharedController}
}
