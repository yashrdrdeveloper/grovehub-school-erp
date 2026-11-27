package launcher

import (
 "context"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry"
 "errors"
 "fmt"
 "net/http"
 "os"
 "os/signal"
 "syscall"
 "sync/atomic"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/readiness"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
)

// Run starts the private application listener with bounded readiness checks.
func Run() error {
 ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
 defer stop()
 db,err:=postgres.OpenProjectPool(ctx,os.Getenv("SCHOOL_DATABASE_URL"))
 if err!=nil {return fmt.Errorf("SCHOOL_DATABASE_URL: %w",err)}
 defer db.Close()
 if err:=registry.Check(ctx,db);err!=nil{return err}
 addr := os.Getenv("GROOVEHUB_HTTP_ADDR")
 if addr == "" { addr = "127.0.0.1:8081" }
 mux := http.NewServeMux()
 if err := registerIdentity(mux); err != nil { return err }
 if err:=registry.RegisterHTTP(mux,db,os.Getenv("SCHOOL_REGISTRY_READ_GRANTS"),os.Getenv("SCHOOL_REGISTRY_RENAME_KEYS"));err!=nil{return err}
 mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
 mux.HandleFunc("GET /health/database",func(w http.ResponseWriter,r *http.Request){
  probe,cancel:=context.WithTimeout(r.Context(),2*time.Second)
  defer cancel()
  if err:=db.PingContext(probe);err!=nil {http.Error(w,"database unavailable",http.StatusServiceUnavailable);return}
  w.WriteHeader(http.StatusNoContent)
 })
 var draining atomic.Bool
 checks:=[]readiness.Check{db.PingContext,func(ctx context.Context)error{return registry.Check(ctx,db)}}
 upstream,err:=readiness.HTTP(os.Getenv("PLATFORM_ORIGIN"),"/ready")
 if err!=nil{return err}
 checks=append(checks,upstream)
 mux.HandleFunc("GET /health/ready",readiness.Handler(&draining,checks...))
 server := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5*time.Second, ReadTimeout: 15*time.Second, WriteTimeout: 30*time.Second, IdleTimeout: 60*time.Second, MaxHeaderBytes: 1<<20}
 done := make(chan error, 1)
 go func() { done <- server.ListenAndServe() }()
 select {
 case err := <-done:
  if errors.Is(err, http.ErrServerClosed) { return nil }
  return err
 case <-ctx.Done():
  draining.Store(true)
  shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  return server.Shutdown(shutdown)
 }
}
