package bootstrap

import (
	"grovehub/cmd/internal/src/command"
	"grovehub/cmd/internal/src/constants"
	databaseimpl "grovehub/cmd/internal/src/operations/database/impl"
	"grovehub/cmd/internal/src/operations/migration"
	"grovehub/cmd/internal/src/registry"
	"grovehub/cmd/internal/src/strategies/auditreplay"
	"grovehub/cmd/internal/src/strategies/auditstorage"
	"grovehub/cmd/internal/src/strategies/bootstrapadmin"
	"grovehub/cmd/internal/src/strategies/mailreconcile"
	"grovehub/cmd/internal/src/strategies/migrate"
	"grovehub/cmd/internal/src/strategies/rbacmigrate"
	"grovehub/cmd/internal/src/strategies/serve"
	"grovehub/services/audit"
	"grovehub/services/foundation"
	"grovehub/services/identity"
	"grovehub/services/infrastructure"
	"grovehub/services/integration"
)

func New() registry.Registry {
	db := databaseimpl.New()
	migrations := map[string]migration.Migration{
		"foundation": migration.Operation(foundation.Migrate),
		"infrastructure": migration.Operation(infrastructure.Migrate),
		"integration": migration.Operation(integration.Migrate),
		"identity": migration.Operation(identity.Migrate),
		"audit": migration.Operation(audit.Migrate),
	}

	return registry.New(map[string]command.Command{
		constants.Serve: serve.New(),
		constants.Migrate: migrate.New(db, migrations),
		constants.AuditReplay: auditreplay.New(db),
		constants.AuditStorage: auditstorage.New(db),
		constants.BootstrapAdmin: bootstrapadmin.New(),
		constants.MailReconcile: mailreconcile.New(db),
		constants.RBACMigrate: rbacmigrate.New(),
	})
}
