// Package bootstrap defines Integration's composition boundary.
package bootstrap

import (
    "net/http"

    integrationclient "grovehub/services/integration/client"
)

// Runtime owns the Integration HTTP controller, internal client and workers.
type Runtime interface {
    RegisterHTTP(*http.ServeMux)
    Client() integrationclient.Client
    LocalOnly() bool
    Close()
}

type Bootstrap interface {
    Build() (Runtime,error)
}
