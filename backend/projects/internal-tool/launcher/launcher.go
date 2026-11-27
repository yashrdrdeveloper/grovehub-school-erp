package launcher

import (
 "context"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/correlation"
 "errors"
 "fmt"
 "net/http"
 "os"
 "os/signal"
 "syscall"
 "sync/atomic"
 auditclient "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/client/audit"
 auditcontroller "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/audit"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/auditdelivery"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/readiness"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
 auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/authentication"
 controller "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/authentication"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/resources/migrations"
 workflow "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
 platform "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/client/platform"
 workflowController "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/workflow"
)

// Run starts the private application listener with bounded readiness checks.
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
 bindings,err:=platform.Bindings(platform.Config{
  EnvironmentCode:os.Getenv("INTERNAL_TOOL_CATALOG_ENVIRONMENT"),
  TemplateID:os.Getenv("INTERNAL_TOOL_CATALOG_TEMPLATE_ID"),
  Origin:os.Getenv("INTERNAL_TOOL_PLATFORM_ORIGIN"),
  KeyID:os.Getenv("INTERNAL_TOOL_WORKLOAD_KEY_ID"),
  PrivateKey:os.Getenv("INTERNAL_TOOL_WORKLOAD_PRIVATE_KEY"),
 })
 if err!=nil{return err}
 workflowController.New(workflow.NewWithBindings(db,bindings),authController).Register(mux)
 publisher,err:=auditclient.New(auditclient.Config{Origin:os.Getenv("INTERNAL_TOOL_AUDIT_ORIGIN"),KeyID:os.Getenv("INTERNAL_TOOL_AUDIT_KEY_ID"),PrivateKey:os.Getenv("INTERNAL_TOOL_AUDIT_PRIVATE_KEY")})
 if err!=nil{return err}
 reader,err:=auditclient.NewReader(auditclient.Config{Origin:os.Getenv("INTERNAL_TOOL_AUDIT_READ_ORIGIN"),KeyID:os.Getenv("INTERNAL_TOOL_AUDIT_READ_KEY_ID"),PrivateKey:os.Getenv("INTERNAL_TOOL_AUDIT_READ_PRIVATE_KEY")})
 if err!=nil{return err}
 if reader!=nil&&publisher!=nil&&(os.Getenv("INTERNAL_TOOL_AUDIT_READ_KEY_ID")==os.Getenv("INTERNAL_TOOL_AUDIT_KEY_ID")||os.Getenv("INTERNAL_TOOL_AUDIT_READ_PRIVATE_KEY")==os.Getenv("INTERNAL_TOOL_AUDIT_PRIVATE_KEY")){return errors.New("audit reader and writer must use independent keys")}
 monitor:=auditdelivery.NewMonitor(db,publisher!=nil)
 auditcontroller.Register(mux,authController,monitor,reader)
 mux.HandleFunc("GET /health/audit-delivery",monitor.HealthHandler)
 if publisher!=nil{
  workerCtx,workerStop:=context.WithCancel(ctx)
  workerDone:=make(chan struct{})
  go func(){defer close(workerDone);auditdelivery.New(db,publisher).Run(workerCtx)}()
  defer func(){workerStop();<-workerDone}()
 }
 mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
 mux.HandleFunc("GET /health/database",func(w http.ResponseWriter,r *http.Request){
  probe,cancel:=context.WithTimeout(r.Context(),2*time.Second)
  defer cancel()
  if err:=db.PingContext(probe);err!=nil {http.Error(w,"database unavailable",http.StatusServiceUnavailable);return}
  w.WriteHeader(http.StatusNoContent)
 })
 var draining atomic.Bool
 checks:=[]readiness.Check{db.PingContext}
 checks=append(checks,func(ctx context.Context)error{
  if err:=migrations.Check(ctx,db);err!=nil{return err}
  var active bool
  if err:=db.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM internal_tool.operators WHERE persona='L2' AND status='ACTIVE')").Scan(&active);err!=nil{return err}
  if !active{return errors.New("no active L2 operator")}
  return nil
 })
 if len(bindings)>0{
  upstream,err:=readiness.HTTP(os.Getenv("INTERNAL_TOOL_PLATFORM_ORIGIN"),"/ready")
  if err!=nil{return err}
  checks=append(checks,upstream)
 }
 mux.HandleFunc("GET /health/ready",readiness.Handler(&draining,checks...))
 server := &http.Server{Addr: addr, Handler: correlation.Ingress(mux), ReadHeaderTimeout: 5*time.Second, ReadTimeout: 15*time.Second, WriteTimeout: 30*time.Second, IdleTimeout: 60*time.Second, MaxHeaderBytes: 1<<20}
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
