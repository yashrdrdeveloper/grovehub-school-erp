package impl

import (
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
	return authpublicimpl.New(services.auth),
		authprotectedimpl.New(services.auth),
		accountpublicimpl.New(services.account),
		accountprotectedimpl.New(services.account)
}
