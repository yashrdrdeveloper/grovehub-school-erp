package impl

import (
	"context"
	"fmt"
	"grovehub/launcher/internal/src/server"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type httpServer struct {
	address  string
	handler  http.Handler
	logger   *slog.Logger
	draining *atomic.Bool
}

var _ server.Server = (*httpServer)(nil)

func New(address string, handler http.Handler, logger *slog.Logger, draining *atomic.Bool) server.Server {
	return &httpServer{address: address, handler: handler, logger: logger, draining: draining}
}

func (s *httpServer) Run(ctx context.Context) error {
	transport := &http.Server{
		Addr:              s.address,
		Handler:           s.handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("bind HTTP listener %s: %w", s.address, err)
	}

	return RunLifecycle(ctx, listener, transport, s.logger, s.draining)
}
