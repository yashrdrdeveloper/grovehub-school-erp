package helpers

import (
	"context"

	"gorm.io/gorm"

	"grovehub/services/audit/internal/app/services/storage"
	"grovehub/services/audit/internal/app/strategies/command"
	"grovehub/services/audit/internal/app/strategies/command/auditstorage/constants"
)

type Inspector func(context.Context, *gorm.DB, string) ([]storage.StorageTable, error)

func Execute(ctx context.Context, database command.Database[*gorm.DB], inspect Inspector, expected string) ([]storage.StorageTable, error) {
	db, closeDB, err := database.Audit(ctx, nil)
	if err != nil {
		return nil, constants.ErrConnection
	}
	defer closeDB()

	tables, err := inspect(ctx, db, expected)
	if err != nil {
		return nil, constants.ErrInspection
	}

	return tables, nil
}
