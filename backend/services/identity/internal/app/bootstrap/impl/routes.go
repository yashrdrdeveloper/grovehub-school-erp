package impl

import (
	cookieimpl "grovehub/services/identity/internal/app/transport/cookies/impl"
	accountprotected "grovehub/services/identity/internal/app/controller/account/protected"
	accountprotectedimpl "grovehub/services/identity/internal/app/controller/account/protected/impl"
	accountpublic "grovehub/services/identity/internal/app/controller/account/public"
	accountpublicimpl "grovehub/services/identity/internal/app/controller/account/public/impl"
	authprotected "grovehub/services/identity/internal/app/controller/authentication/protected"
	authprotectedimpl "grovehub/services/identity/internal/app/controller/authentication/protected/impl"
	authpublic "grovehub/services/identity/internal/app/controller/authentication/public"
	authpublicimpl "grovehub/services/identity/internal/app/controller/authentication/public/impl"
)

func (b *builder) newControllers(services serviceGraph) (authpublic.AuthenticationController, authprotected.AuthenticationController, accountpublic.AccountController, accountprotected.AccountController) {
	manager := cookieimpl.New()
	return authpublicimpl.New(services.auth, manager),
		authprotectedimpl.New(services.auth, manager),
		accountpublicimpl.New(services.account, manager),
		accountprotectedimpl.New(services.account)
}
