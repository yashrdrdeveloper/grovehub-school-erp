package launcher

import (
 "context"
 "io"
 "log/slog"
 "net/http"
 coregateway "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 "grovehub/launcher/session"
 external "grovehub/launcher/internal/src/services/gateway/operations/external/impl"
 "grovehub/launcher/internal/src/services/gateway/operations/workload"
 tenantaccess "grovehub/launcher/internal/src/services/gateway/operations/tenant/access"
 "gorm.io/gorm"
 bootstrap "grovehub/launcher/internal/src/bootstrap/impl"
)

func Run(logger *slog.Logger) error {
 runtime,err:=bootstrap.New(logger).Build()
 if err!=nil{return err}
 defer runtime.Close()
 return runtime.Run()
}

// Operator commands share the same Platform-owned database and logger setup.
func NewLogger(output io.Writer)(*slog.Logger,error){return bootstrap.NewLogger(output)}
func NewPostgresDB(ctx context.Context)(*gorm.DB,error){return bootstrap.NewPostgresDB(ctx)}
func NewAuditPostgresDB(ctx context.Context,shared *gorm.DB)(*gorm.DB,func(),error){return bootstrap.NewAuditPostgresDB(ctx,shared)}
var ErrDatabaseURLRequired=bootstrap.ErrDatabaseURLRequired

 // DelegatedIdentity is the verified actor boundary for composing the gateway.
 type DelegatedIdentity interface {
   VerifyDelegatedSession(context.Context,string)(session.Principal,error)
   session.ProtectedRouteAuthorizer
   session.AccessVerifier
 }
 type TenantAccessConfig = tenantaccess.TenantAccessConfig
 func NewTenantState(c TenantAccessConfig)(session.TenantStateLookup,error){return tenantaccess.NewTenantState(c)}
 func NewWorkloadVerifier(raw string,identity DelegatedIdentity,webhook coregateway.WorkloadVerifier,checks ...session.TenantCheck)(coregateway.WorkloadVerifier,error){
  return workload.NewPlatformWorkloads(raw,identity,webhook,checks...)
 }
 func NewExternalGateway(mux *http.ServeMux,identity session.AccessVerifier,workloads coregateway.WorkloadVerifier,origin string,limits ...map[string]int)http.Handler{
  return external.New(mux,identity,workloads,origin,limits...)
 }
