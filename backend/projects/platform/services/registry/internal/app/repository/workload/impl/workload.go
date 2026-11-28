package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/constants"
	"grovehub/services/registry/internal/app/repository/workload"
)

type WorkloadRepository struct {
	db *gorm.DB
}

func NewWorkloadRepository(db *gorm.DB) workload.WorkloadRepository {
	return &WorkloadRepository{db: db}
}

func (r *WorkloadRepository) Find(ctx context.Context, workloadID string) (entity.Workload, error) {
	var result entity.Workload
	err := coredatabase.TransactionDB(ctx, r.db).Where("workload_id = ?", workloadID).Take(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Workload{}, constants.ErrNotFound
	}
	if err != nil {
		return entity.Workload{}, err
	}

	return result, nil
}
