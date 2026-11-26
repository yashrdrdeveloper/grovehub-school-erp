package impl

import (
	"net/http"
	roleprotected "grovehub/services/identity/internal/app/controller/role/protected"
	roleprotectedimpl "grovehub/services/identity/internal/app/controller/role/protected/impl"
	accountroleprotected "grovehub/services/identity/internal/app/controller/accountrole/protected"
	accountroleprotectedimpl "grovehub/services/identity/internal/app/controller/accountrole/protected/impl"
	approvalprotected "grovehub/services/identity/internal/app/controller/approval/protected"
	approvalprotectedimpl "grovehub/services/identity/internal/app/controller/approval/protected/impl"

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
	roleController roleprotected.RoleController
	accountroleController accountroleprotected.AccountRoleController
	approvalController approvalprotected.ApprovalController
}

func (b *builder) newControllers(services serviceGraph) controllerGraph {
	manager := cookieimpl.New()
	return controllerGraph{
		publicAuth:       authpublicimpl.New(services.auth, manager),
		protectedAuth:    authprotectedimpl.New(services.auth, manager),
		publicActivation: activationpublicimpl.New(services.account, manager),
		protectedAccount: accountprotectedimpl.New(services.account, services.catalog),
		roleController: roleprotectedimpl.New(services.roles),
		accountroleController: accountroleprotectedimpl.New(services.accountRoles),
		approvalController: approvalprotectedimpl.NewWithMutations(services.approvalReads, services.approvalMutations, nil),
	}
}

func (controllers controllerGraph) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, controllers.publicAuth, controllers.protectedAuth, controllers.publicActivation, controllers.protectedAccount, controllers.roleController, controllers.accountroleController, controllers.approvalController)
}

func (module *runtime) RouteLimits() map[string]int {
	return routes.RouteLimits()
}
