package server

import (
	"context"
	"net"
)

// @Mock
type Lifecycle interface {
	Serve(net.Listener) error
	Shutdown(context.Context) error
	Close() error
}
