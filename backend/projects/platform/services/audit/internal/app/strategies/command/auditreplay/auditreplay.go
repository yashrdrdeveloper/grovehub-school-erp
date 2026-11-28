package auditreplay

import (
	"gorm.io/gorm"

	"grovehub/services/audit/internal/app/strategies/command"
	"grovehub/services/audit/internal/app/strategies/command/auditreplay/helpers"
	"grovehub/services/audit/internal/app/strategies/command/auditreplay/impl"
)

type ReplayOperation = helpers.ReplayOperation

func New(db command.Database[*gorm.DB]) command.Command {
	return impl.New(db)
}

func NewWithReplay(db command.Database[*gorm.DB], operation ReplayOperation) command.Command {
	return impl.NewWithReplay(db, operation)
}
