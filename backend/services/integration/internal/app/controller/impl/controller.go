package impl

import controller "grovehub/services/integration/internal/app/controller"

// Controller composes per-resource HTTP handlers without changing registered URLs.
type Controller struct {
    controller.ConnectionController
    controller.ProviderController
}

var _ controller.Controller = (*Controller)(nil)

func New(connections controller.ConnectionController, providers controller.ProviderController) controller.Controller {
    return &Controller{ConnectionController:connections,ProviderController:providers}
}
