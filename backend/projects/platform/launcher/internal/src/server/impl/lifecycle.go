package impl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	coretelemetry "grovehub/libs/core/telemetry"
	"grovehub/launcher/internal/src/server"
)

// RunLifecycle owns serving and graceful shutdown for a bound listener.
func RunLifecycle(ctx context.Context, listener net.Listener, transport server.Lifecycle, logger *slog.Logger, draining *atomic.Bool) error {
	edgeService := coretelemetry.ServiceLauncherEdge
	logger.Info("HTTP server started", "service", edgeService.Name, "serviceKey", edgeService.Key, "address", listener.Addr().String())
	done := make(chan error, 1)
	go func() {
		done <- transport.Serve(listener)
	}()
	select {
	case err := <-done:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	case <-ctx.Done():
		draining.Store(true)
		shutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := transport.Shutdown(shutdown); err != nil {
			_ = transport.Close()
			return err
		}
		return nil
	}
}
