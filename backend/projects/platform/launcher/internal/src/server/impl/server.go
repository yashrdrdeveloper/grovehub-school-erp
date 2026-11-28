package impl

import (
 "context";"errors";"fmt";"log/slog";"net";"net/http";"sync/atomic";"time"
 "grovehub/launcher/internal/src/server"
 coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
)
type httpServer struct {address string;handler http.Handler;logger *slog.Logger;draining *atomic.Bool}
var _ server.Server = (*httpServer)(nil)
func New(address string,handler http.Handler,logger *slog.Logger,draining *atomic.Bool) server.Server {
 return &httpServer{address:address,handler:handler,logger:logger,draining:draining}
}
func(s *httpServer) Run(ctx context.Context)error{
 edgeService:=coretelemetry.ServiceLauncherEdge
	server := &http.Server{
		Addr: s.address,
		Handler: s.handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 45 * time.Second,
		IdleTimeout: 60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("bind HTTP listener %s: %w", s.address, err)
	}

	s.logger.Info(
		"HTTP server started",
		"service", edgeService.Name,
		"serviceKey", edgeService.Key,
		"address", listener.Addr().String(),
	)

	done:=make(chan error,1)
	go func(){done<-server.Serve(listener)}()
	select {
	case err:=<-done:
		if err!=nil&&!errors.Is(err,http.ErrServerClosed){return fmt.Errorf("serve HTTP: %w",err)}
		return nil
	case <-ctx.Done():
		s.draining.Store(true)
		shutdown,cancel:=context.WithTimeout(context.Background(),15*time.Second)
		defer cancel()
		if err:=server.Shutdown(shutdown);err!=nil{_ = server.Close();return err}
		return nil
	}
}
