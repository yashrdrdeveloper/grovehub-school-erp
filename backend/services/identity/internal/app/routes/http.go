package routes

import (
	"net/http"

	accountprotected "grovehub/services/identity/internal/app/controller/account/protected"
	activationpublic "grovehub/services/identity/internal/app/controller/activation/public"
	authprotected "grovehub/services/identity/internal/app/controller/authentication/protected"
	authpublic "grovehub/services/identity/internal/app/controller/authentication/public"
)

// RegisterHTTP only composes annotated, generated route groups.
// Authentication and authorization are enforced by the Launcher gateway.
// @Routes /identity
func RegisterHTTP(mux *http.ServeMux, publicAuth authpublic.AuthenticationController, protectedAuth authprotected.AuthenticationController, publicActivation activationpublic.ActivationController, protectedAccount accountprotected.AccountController) {
	registerPublicAuthenticationRoutes(mux, publicAuth)
	registerProtectedAuthenticationRoutes(mux, protectedAuth)
	registerPublicActivationRoutes(mux, publicActivation)
	registerProtectedAccountRoutes(mux, protectedAccount)
}
