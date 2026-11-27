package impl

import (
	"context"
	"time"

	"gorm.io/gorm"
	"grovehub/services/audit/internal/app/services/monitoring"
	"grovehub/services/infrastructure/client/messaging"
)

func (module *runtime) Ready() bool {
	return module.monitor.Ready()
}

func (module *runtime) probe(ctx context.Context,db *gorm.DB,client messaging.MessagingClient) {
	defer module.wait.Done()
	ticker:=time.NewTicker(monitoring.ProbeInterval)
	defer ticker.Stop()
	for ctx.Err()==nil {
		bounded,cancel:=context.WithTimeout(ctx,monitoring.ProbeTimeout)
		sqlDB,err:=db.DB()
		ready:=err==nil && sqlDB.PingContext(bounded)==nil && client.Health(bounded)==nil
		cancel()
		module.monitor.Probe(time.Now(),ready)
		select {
		case <-ctx.Done(): return
		case <-ticker.C:
		}
	}
}
