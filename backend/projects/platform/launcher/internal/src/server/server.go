package server

import (
 "context"
 "net"
)

type Server interface { Run(context.Context) error }

// Lifecycle is the HTTP transport lifecycle managed by the server operation.
type Lifecycle interface {
 Serve(net.Listener) error
 Shutdown(context.Context) error
 Close() error
}
