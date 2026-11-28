package auditstorage

import (
	"gorm.io/gorm"

	"grovehub/services/audit/internal/app/strategies/command"
	"grovehub/services/audit/internal/app/strategies/command/auditstorage/helpers"
	"grovehub/services/audit/internal/app/strategies/command/auditstorage/impl"
)

type Inspector = helpers.Inspector

func New(db command.Database[*gorm.DB]) command.Command {
	return impl.New(db)
}

func NewWithInspector(db command.Database[*gorm.DB], operation Inspector) command.Command {
	return impl.NewWithInspector(db, operation)
}
