// Package bootstrap defines Foundation's composition contract without
// exposing repositories, controllers or event transport internals.
package bootstrap

import (
    coremodule "grovehub/libs/core/module"
    foundationclient "grovehub/services/foundation/client/person"
)

// Runtime is Foundation's public module contract: common lifecycle plus its
// typed in-process Person client.
type Runtime interface {
    coremodule.Module
    Client() foundationclient.Client
}

// Bootstrap assembles all Foundation dependencies and starts owned workers.
type Bootstrap interface {
    Build() (Runtime, error)
}
