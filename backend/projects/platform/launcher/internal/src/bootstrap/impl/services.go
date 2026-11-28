package impl

import (
 "context"
 "fmt"
 "time"

 "gorm.io/gorm"
 "grovehub/launcher/internal/src/bootstrap/helpers"
 "grovehub/internal/session"
 "grovehub/services/audit"
 "grovehub/services/foundation"
 "grovehub/services/identity"
 "grovehub/services/infrastructure"
 messaging "grovehub/services/infrastructure/client/messaging"
 "grovehub/services/integration"
)

type serviceModules struct {
 infrastructure infrastructure.Module
 foundation foundation.Module
 integration integration.Module
 identity identity.Module
 audit audit.Module
 auditDB *gorm.DB
}

func(b *builder) newServices(ctx context.Context,db *gorm.DB,tenantState session.TenantStateLookup,cleanup *helpers.Cleanup)(serviceModules,error){
 var modules serviceModules
 var err error
 modules.infrastructure,err=infrastructure.New(messaging.DefaultMessaging(),db,b.logger)
 if err!=nil{return serviceModules{},fmt.Errorf("initialize Infrastructure: %w",err)}
 cleanup.Add(func(){_=modules.infrastructure.Close()})

 modules.foundation,err=foundation.New(b.logger,db,modules.infrastructure.EventStore())
 if err!=nil{return serviceModules{},fmt.Errorf("initialize Foundation: %w",err)}
 cleanup.Add(modules.foundation.Close)

 modules.integration,err=integration.NewWithTenantState(b.logger,db,modules.infrastructure.MessagingClient(),tenantState)
 if err!=nil{return serviceModules{},fmt.Errorf("initialize Integration: %w",err)}
 cleanup.Add(modules.integration.Close)

 modules.identity,err=identity.New(b.logger,db,modules.foundation.Client(),modules.integration.Client(),modules.infrastructure.EventStore(),tenantState)
 if err!=nil{return serviceModules{},fmt.Errorf("initialize Identity: %w",err)}
 cleanup.Add(modules.identity.Close)

 auditCtx,cancelAudit:=context.WithTimeout(ctx,5*time.Second)
 auditDB,closeAuditDB,err:=NewAuditPostgresDB(auditCtx,db)
 cancelAudit()
 if err!=nil{return serviceModules{},err}
 cleanup.Add(closeAuditDB)
 modules.auditDB=auditDB
 modules.audit,err=audit.New(b.logger,auditDB,modules.infrastructure.MessagingClient(),modules.identity)
 if err!=nil{return serviceModules{},fmt.Errorf("initialize Audit: %w",err)}
 cleanup.Add(modules.audit.Close)
 return modules,nil
}
