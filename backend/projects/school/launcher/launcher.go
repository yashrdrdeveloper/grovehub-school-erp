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
 corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
 "time"

 coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

// Run starts the private application listener with bounded readiness checks.
func Run() error {
 ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
 defer stop()
 db,err:=coredatabase.OpenProjectPool(ctx,os.Getenv("SCHOOL_DATABASE_URL"))
 if err!=nil {return fmt.Errorf("SCHOOL_DATABASE_URL: %w",err)}
 defer db.Close()
 if err:=registry.Check(ctx,db);err!=nil{return err}
 addr := os.Getenv("GROOVEHUB_HTTP_ADDR")
 if addr == "" { addr = "127.0.0.1:8081" }
 mux := http.NewServeMux()
 if err := registerIdentity(mux); err != nil { return err }
 if err:=registry.RegisterHTTP(mux,db,os.Getenv("SCHOOL_REGISTRY_READ_GRANTS"),os.Getenv("SCHOOL_REGISTRY_RENAME_KEYS"),os.Getenv("SCHOOL_REGISTRY_STATUS_KEYS"),os.Getenv("SCHOOL_REGISTRY_ACCESS_KEYS"));err!=nil{return err}
 deliveryConfig:=registry.DeliveryConfig{Origin:os.Getenv("SCHOOL_RECEIPT_ORIGIN"),SchoolID:os.Getenv("SCHOOL_RECEIPT_SCHOOL_ID"),KeyID:os.Getenv("SCHOOL_RECEIPT_KEY_ID"),PrivateKey:os.Getenv("SCHOOL_RECEIPT_PRIVATE_KEY")}
 configs,err:=registry.ConfigureDelivery(os.Getenv("SCHOOL_RECEIPT_PUBLISHERS"),deliveryConfig,os.Getenv("SCHOOL_WORKLOAD_KEY_ID"),os.Getenv("SCHOOL_WORKLOAD_PRIVATE_KEY"),os.Getenv("SCHOOL_REGISTRY_READ_GRANTS"),os.Getenv("SCHOOL_REGISTRY_RENAME_KEYS"),os.Getenv("SCHOOL_REGISTRY_STATUS_KEYS"),os.Getenv("SCHOOL_REGISTRY_ACCESS_KEYS"))
 if err!=nil{return err}
 stopDelivery,err:=registry.StartDeliveries(ctx,db,configs);if err!=nil{return err};defer stopDelivery()
 mux.HandleFunc("GET /health/audit-delivery",registry.DeliveryHealth(db,configs))
 mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
 mux.HandleFunc("GET /health/database",func(w http.ResponseWriter,r *http.Request){
  probe,cancel:=context.WithTimeout(r.Context(),2*time.Second)
  defer cancel()
  if err:=db.PingContext(probe);err!=nil {http.Error(w,"database unavailable",http.StatusServiceUnavailable);return}
  w.WriteHeader(http.StatusNoContent)
 })
 var draining atomic.Bool
 checks:=[]corehttp.Check{db.PingContext,func(ctx context.Context)error{return registry.Check(ctx,db)}}
 upstream,err:=corehttp.ReadinessHTTP(os.Getenv("PLATFORM_ORIGIN"),"/ready")
 if err!=nil{return err}
 checks=append(checks,upstream)
 mux.HandleFunc("GET /health/ready",corehttp.ReadinessHandler(&draining,checks...))
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
