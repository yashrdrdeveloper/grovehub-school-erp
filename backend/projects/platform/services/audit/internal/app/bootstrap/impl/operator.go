package impl

import (
	"context"

	"gorm.io/gorm"

	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/libs/platform/command"
	"grovehub/services/audit/internal/app/bootstrap/constants"
	"grovehub/services/audit/internal/app/strategies/command/auditreplay"
	"grovehub/services/audit/internal/app/strategies/command/auditstorage"
)

func NewOperatorBootstrap(migrate func(context.Context, *gorm.DB) error) platformbootstrap.Service[*gorm.DB] {
	return platformbootstrap.Registration[*gorm.DB]{
		ServiceName: constants.ServiceName,
		RegisterCommands: func(db command.Database[*gorm.DB]) map[string]command.Command {
			return map[string]command.Command{
				constants.AuditReplay:  auditreplay.New(db),
				constants.AuditStorage: auditstorage.New(db),
			}
		},
		MigrationPlan: platformbootstrap.Migration[*gorm.DB]{
			Target: platformbootstrap.AuditDatabase,
			Run:    migrate,
		},
	}
}
