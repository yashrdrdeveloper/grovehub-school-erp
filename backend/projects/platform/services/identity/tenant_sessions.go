package identity

import (
 "context"
 "log/slog"

 "gorm.io/gorm"
 schoolgateway "grovehub/internal/session"
 sessionrepo "grovehub/services/identity/internal/app/repository/session/impl"
 "grovehub/services/identity/internal/app/services/authentication/operations/reconciletenant"
)

// StartTenantSessionReconciler owns the background cleanup lifecycle. The
// returned stop function waits for the worker before the caller closes DB.
func StartTenantSessionReconciler(parent context.Context, logger *slog.Logger, db *gorm.DB, lookup schoolgateway.TenantStateLookup) func() {
 if db == nil || lookup == nil { return func() {} }
 ctx, cancel := context.WithCancel(parent)
 done := make(chan struct{})
 worker := reconciletenant.New(sessionrepo.NewTenantSessions(db), lookup, logger)
 go func() {
  defer close(done)
  worker.Run(ctx)
 }()
 return func() { cancel(); <-done }
}
