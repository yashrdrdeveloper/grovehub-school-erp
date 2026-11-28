package audit

import (
	"context"

	"gorm.io/gorm"

	"grovehub/services/audit/internal/app/services/storage"
)

type StorageTable = storage.StorageTable

func InspectStorage(ctx context.Context, db *gorm.DB, expectedDatabase string) ([]StorageTable, error) {
	return storage.InspectStorage(ctx, db, expectedDatabase)
}
