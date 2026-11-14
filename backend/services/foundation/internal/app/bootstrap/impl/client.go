package impl

import (
    "net/http"

    foundationclient "grovehub/services/foundation/client/person"
    clientimpl "grovehub/services/foundation/internal/app/client/person/impl"
    protected "grovehub/services/foundation/internal/app/controller/person/protected"
    "grovehub/services/foundation/internal/app/routes"
)

// The in-process client uses the same generated HTTP routes and controller
// validation pipeline as public requests; it does not open a TCP connection.
func (bootstrap *builder) newClient(controller protected.PersonController) foundationclient.Client {
    internalMux := http.NewServeMux()
    routes.RegisterHTTP(internalMux, controller)
    return clientimpl.NewPersonClient(internalMux)
}
