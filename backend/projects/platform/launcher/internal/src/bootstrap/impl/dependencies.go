package impl

import (
	"grovehub/launcher/internal/src/bootstrap"
	serverimpl "grovehub/launcher/internal/src/server/impl"
	"grovehub/services/audit"
	"grovehub/services/foundation"
	"grovehub/services/identity"
	"grovehub/services/infrastructure"
	"grovehub/services/integration"
)

// DefaultDependencies is the production composition root. Each adapter has
// explicit ownership; the builder remains responsible for rollback.
func DefaultDependencies() bootstrap.Dependencies {
	return bootstrap.Dependencies{
		OpenDatabase:          NewPostgresDB,
		OpenAuditDatabase:     NewAuditPostgresDB,
		NewInfrastructure:     infrastructure.New,
		NewFoundation:         foundation.New,
		NewIntegration:        integration.NewWithTenantState,
		NewIdentity:           identity.New,
		NewAudit:              audit.New,
		StartTenantReconciler: identity.StartTenantSessionReconciler,
		NewServer:             serverimpl.New,
	}
}
