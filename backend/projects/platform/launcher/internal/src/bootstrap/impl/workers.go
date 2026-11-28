package impl

import (
 "context"
 "gorm.io/gorm"
 "grovehub/launcher/internal/src/bootstrap/helpers"
 "grovehub/internal/session"
)

func(b *builder) startWorkers(ctx context.Context,db *gorm.DB,lookup session.TenantStateLookup,cleanup *helpers.Cleanup){
 if lookup!=nil{cleanup.Add(b.dependencies.StartTenantReconciler(ctx,b.logger,db,lookup))}
}
