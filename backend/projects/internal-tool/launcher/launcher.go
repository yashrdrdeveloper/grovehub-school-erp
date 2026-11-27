package launcher

import (
 "context"
 "errors"
 "net/http"
 "os"
 "os/signal"
 "syscall"
 "time"
)

// Run starts a private migration-stage listener. Readiness remains false
// until application authentication and dependencies are configured.
func Run() error {
 ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
 defer stop()
 addr := os.Getenv("GROOVEHUB_HTTP_ADDR")
 if addr == "" { addr = "127.0.0.1:8082" }
 mux := http.NewServeMux()
 mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
 mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "application dependencies not initialized", http.StatusServiceUnavailable) })
 server := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5*time.Second, ReadTimeout: 15*time.Second, WriteTimeout: 30*time.Second, IdleTimeout: 60*time.Second, MaxHeaderBytes: 1<<20}
 done := make(chan error, 1)
 go func() { done <- server.ListenAndServe() }()
 select {
 case err := <-done:
  if errors.Is(err, http.ErrServerClosed) { return nil }
  return err
 case <-ctx.Done():
  shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  return server.Shutdown(shutdown)
 }
}
