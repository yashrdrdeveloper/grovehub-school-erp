package routes

import (
    "net/http"

    protected "grovehub/services/foundation/controller/person/protected"
)

// RegisterHTTP composes Foundation's generated, traced route registrations.
// @Routes /foundation
func RegisterHTTP(mux *http.ServeMux, personController protected.PersonController) {
    registerProtectedPersonRoutes(mux, personController)
}
