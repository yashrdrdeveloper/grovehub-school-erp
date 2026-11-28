package impl

import (
	"context"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/helpers"
	"grovehub/services/registry/internal/app/repository/grant"
	queries "grovehub/services/registry/internal/app/repository/grant/constants"
)

type GrantRepository struct {
	db *gorm.DB
}

func NewGrantRepository(db *gorm.DB) grant.GrantRepository {
	return &GrantRepository{db: db}
}

func (r *GrantRepository) Create(ctx context.Context, value entity.WorkloadGrant) error {
	return helpers.WriteError(coredatabase.TransactionDB(ctx, r.db).Create(&value).Error)
}

func (r *GrantRepository) Revoke(ctx context.Context, workloadID, grantID string) error {
	return helpers.UpdateError(coredatabase.TransactionDB(ctx, r.db).Exec(queries.Revoke, workloadID, grantID))
}
