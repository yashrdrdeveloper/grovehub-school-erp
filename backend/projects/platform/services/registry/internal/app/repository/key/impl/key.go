package impl

import (
	"context"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/helpers"
	"grovehub/services/registry/internal/app/repository/key"
	queries "grovehub/services/registry/internal/app/repository/key/constants"
)

type KeyRepository struct {
	db *gorm.DB
}

func NewKeyRepository(db *gorm.DB) key.KeyRepository {
	return &KeyRepository{db: db}
}

func (r *KeyRepository) Create(ctx context.Context, value entity.WorkloadKey) error {
	return helpers.WriteError(coredatabase.TransactionDB(ctx, r.db).Create(&value).Error)
}

func (r *KeyRepository) Revoke(ctx context.Context, workloadID, keyID string) error {
	return helpers.UpdateError(coredatabase.TransactionDB(ctx, r.db).Exec(queries.Revoke, workloadID, keyID))
}
