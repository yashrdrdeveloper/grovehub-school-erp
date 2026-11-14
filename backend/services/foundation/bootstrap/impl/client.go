package impl

import (
    "net/http"

    foundationclient "grovehub/services/foundation/client/person"
    clientimpl "grovehub/services/foundation/client/person/impl"
    protected "grovehub/services/foundation/controller/person/protected"
    "grovehub/services/foundation/routes"
)

// The in-process client uses the same generated HTTP routes and controller
// validation pipeline as public requests; it does not open a TCP connection.
func (bootstrap *builder) newClient(controller protected.PersonController) foundationclient.Client {
    internalMux := http.NewServeMux()
    routes.RegisterHTTP(internalMux, controller)
    return clientimpl.NewPersonClient(internalMux)
}
