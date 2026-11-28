package impl

import (
	"context"
	"database/sql"
	"net/http"
	"sync/atomic"
	"time"

	"gorm.io/gorm"

	health "grovehub/launcher/internal/src/controller/impl"
	"grovehub/launcher/internal/src/routes"
	"grovehub/services/identity"
	"grovehub/services/integration"
)

func registerRoutes(mux *http.ServeMux, db *gorm.DB, sqlDB *sql.DB, modules serviceModules, draining *atomic.Bool) error {
	modules.foundation.RegisterHTTP(mux)
	modules.integration.RegisterHTTP(mux)
	integration.RegisterMailReconciliationHTTP(mux, integration.NewMailReconciliationClient(db, identity.NewInvitationDeliveryReader(db)))
	modules.identity.RegisterHTTP(mux)
	modules.audit.RegisterHTTP(mux)
	auditSQL, err := modules.auditDB.DB()
	if err != nil {
		return err
	}

	routes.RegisterHTTP(mux, health.New(func() bool {
		if draining.Load() || !modules.audit.Ready() {
			return false
		}

		probe, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if sqlDB.PingContext(probe) != nil || auditSQL.PingContext(probe) != nil {
			return false
		}

		return modules.infrastructure.MessagingClient().Health(probe) == nil && !draining.Load()
	}))
	routes.RegisterFallback(mux)
	return nil
}
