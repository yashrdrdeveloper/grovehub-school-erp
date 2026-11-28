// Package client preserves the public Integration API while the implementation
// follows Foundation's client/<resource> layout.
package client

import connections "grovehub/services/integration/client/connection"

type Client = connections.Client
