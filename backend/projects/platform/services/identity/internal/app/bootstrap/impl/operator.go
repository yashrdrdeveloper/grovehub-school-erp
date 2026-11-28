package impl

import (
	"context"

	"gorm.io/gorm"

	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/libs/platform/command"
	"grovehub/services/identity/internal/app/strategies/command/bootstrapadmin"
	"grovehub/services/identity/internal/app/strategies/command/mailreconcile"
	"grovehub/services/identity/internal/app/strategies/command/rbacmigrate"
)

func NewOperatorBootstrap(migrate func(context.Context, *gorm.DB) error) platformbootstrap.Service[*gorm.DB] {
	return platformbootstrap.Registration[*gorm.DB]{
		ServiceName: "identity",
		RegisterCommands: func(db command.Database[*gorm.DB]) map[string]command.Command {
			return map[string]command.Command{
				"identity-bootstrap-admin": bootstrapadmin.New(),
				"identity-mail-reconcile":  mailreconcile.New(db),
				"identity-rbac-migrate":    rbacmigrate.New(),
			}
		},
		MigrationPlan: platformbootstrap.Migration[*gorm.DB]{
			Target: platformbootstrap.BusinessDatabase,
			Run:    migrate,
		},
	}
}
