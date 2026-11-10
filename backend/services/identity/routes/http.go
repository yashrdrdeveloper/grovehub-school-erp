package routes

import (
    "net/http"

    "grovehub/services/identity/controller"
)

// RegisterHTTP composes Identity's generated, traced route registrations.
func RegisterHTTP(mux *http.ServeMux, authController controller.AuthController) {
    registerAuthRoutes(mux, authController)
}
