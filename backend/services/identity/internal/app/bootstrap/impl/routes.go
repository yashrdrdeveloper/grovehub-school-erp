package impl

import (
	"net/http"

	accountprotected "grovehub/services/identity/internal/app/controller/account/protected"
	accountprotectedimpl "grovehub/services/identity/internal/app/controller/account/protected/impl"
	activationpublic "grovehub/services/identity/internal/app/controller/activation/public"
	activationpublicimpl "grovehub/services/identity/internal/app/controller/activation/public/impl"
	authprotected "grovehub/services/identity/internal/app/controller/authentication/protected"
	authprotectedimpl "grovehub/services/identity/internal/app/controller/authentication/protected/impl"
	authpublic "grovehub/services/identity/internal/app/controller/authentication/public"
	authpublicimpl "grovehub/services/identity/internal/app/controller/authentication/public/impl"
	"grovehub/services/identity/internal/app/routes"
	cookieimpl "grovehub/services/identity/internal/app/transport/cookies/impl"
)

type controllerGraph struct {
	publicAuth       authpublic.AuthenticationController
	protectedAuth    authprotected.AuthenticationController
	publicActivation activationpublic.ActivationController
	protectedAccount accountprotected.AccountController
}

func (b *builder) newControllers(services serviceGraph) controllerGraph {
	manager := cookieimpl.New()
	return controllerGraph{
		publicAuth:       authpublicimpl.New(services.auth, manager),
		protectedAuth:    authprotectedimpl.New(services.auth, manager),
		publicActivation: activationpublicimpl.New(services.account, manager),
		protectedAccount: accountprotectedimpl.New(services.account),
	}
}

func (controllers controllerGraph) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, controllers.publicAuth, controllers.protectedAuth, controllers.publicActivation, controllers.protectedAccount)
}

func (module *runtime) RouteLimits() map[string]int {
	return routes.RouteLimits()
}
