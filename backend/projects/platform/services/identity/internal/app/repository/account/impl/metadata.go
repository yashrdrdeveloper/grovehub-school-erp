package impl

import (
	"context"
	"math"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/account/constants"
)

type MetadataWriter struct {
	db *gorm.DB
}

func NewMetadataWriter(db *gorm.DB) account.MetadataWriter {
	return &MetadataWriter{db: db}
}

func (writer *MetadataWriter) UpdateMetadata(ctx context.Context, request account.MetadataUpdate) error {
	if writer.db == nil || !coredatabase.TransactionActive(ctx) || request.TenantID == "" || request.AccountID == "" || request.Version < 1 || request.Version == math.MaxInt64 || request.Username == "" || request.UpdatedAt.IsZero() || (request.Status != entity.Active && request.Status != entity.PendingSetup) {
		return account.ErrConflict
	}
	result := coredatabase.TransactionDB(ctx, writer.db).Exec(constants.UpdateMetadata, request.Username, request.SearchDetails, request.UpdatedAt, request.TenantID, request.AccountID, request.Version, request.Status)
	if result.Error != nil {
		return mapWriteError(result.Error)
	}
	if result.RowsAffected != 1 {
		return account.ErrVersionConflict
	}
	return nil
}
