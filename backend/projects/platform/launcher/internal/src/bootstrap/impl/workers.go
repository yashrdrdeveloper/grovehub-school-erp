package impl

import (
 "context"
 "gorm.io/gorm"
 "grovehub/launcher/internal/src/bootstrap/helpers"
 "grovehub/internal/session"
 "grovehub/services/identity"
)

func(b *builder) startWorkers(ctx context.Context,db *gorm.DB,lookup session.TenantStateLookup,cleanup *helpers.Cleanup){
 if lookup!=nil{cleanup.Add(identity.StartTenantSessionReconciler(ctx,b.logger,db,lookup))}
}
