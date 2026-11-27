package launcher

import (
 "context"
 "errors"
 "fmt"
 "net/http"
 "os"
 "os/signal"
 "syscall"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
 auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/authentication"
 controller "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/authentication"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/resources/migrations"
 workflow "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
 workflowController "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/workflow"
)

// Run starts a private migration-stage listener. Readiness remains false
// until application authentication and dependencies are configured.
func Run() error {
 ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
 defer stop()
 db,err:=postgres.OpenProjectPool(ctx,os.Getenv("INTERNAL_TOOL_DATABASE_URL"))
 if err!=nil {return fmt.Errorf("INTERNAL_TOOL_DATABASE_URL: %w",err)}
 defer db.Close()
 addr := os.Getenv("GROOVEHUB_HTTP_ADDR")
 if addr == "" { addr = "127.0.0.1:8082" }
 check,cancel:=context.WithTimeout(ctx,5*time.Second)
 err=migrations.Check(check,db)
 cancel()
 if err!=nil{return err}
 authController,err:=controller.New(auth.New(db),os.Getenv("INTERNAL_TOOL_ALLOWED_ORIGIN"))
 if err!=nil{return err}
 mux := http.NewServeMux()
 authController.Register(mux)
 workflowController.New(workflow.New(db),authController).Register(mux)
 mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
 mux.HandleFunc("GET /health/database",func(w http.ResponseWriter,r *http.Request){
  probe,cancel:=context.WithTimeout(r.Context(),2*time.Second)
  defer cancel()
  if err:=db.PingContext(probe);err!=nil {http.Error(w,"database unavailable",http.StatusServiceUnavailable);return}
  w.WriteHeader(http.StatusNoContent)
 })
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
