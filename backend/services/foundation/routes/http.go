package routes

import (
    "net/http"

    "grovehub/services/foundation/controller"
)

// RegisterHTTP composes Foundation's generated, traced route registrations.
func RegisterHTTP(mux *http.ServeMux, personController controller.PersonController) {
    registerPersonRoutes(mux, personController)
}
