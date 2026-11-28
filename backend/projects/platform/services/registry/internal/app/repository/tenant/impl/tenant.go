package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/constants"
	"grovehub/services/registry/internal/app/repository/tenant"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) tenant.TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Find(ctx context.Context, projectID, tenantID string) (entity.Tenant, error) {
	var result entity.Tenant
	err := coredatabase.TransactionDB(ctx, r.db).Where("project_id = ? AND tenant_id = ?", projectID, tenantID).Take(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Tenant{}, constants.ErrNotFound
	}
	if err != nil {
		return entity.Tenant{}, err
	}

	return result, nil
}
