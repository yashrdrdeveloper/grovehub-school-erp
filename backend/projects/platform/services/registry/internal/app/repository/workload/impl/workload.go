package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/helpers"
	"grovehub/services/registry/internal/app/repository/constants"
	"grovehub/services/registry/internal/app/repository/workload"
	queries "grovehub/services/registry/internal/app/repository/workload/constants"
)

type WorkloadRepository struct {
	db *gorm.DB
}

func NewWorkloadRepository(db *gorm.DB) workload.WorkloadRepository {
	return &WorkloadRepository{db: db}
}

func (r *WorkloadRepository) Find(ctx context.Context, workloadID string) (entity.Workload, error) {
	var result entity.Workload
	err := coredatabase.TransactionDB(ctx, r.db).Where(queries.FindByID, workloadID).Take(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Workload{}, constants.ErrNotFound
	}
	if err != nil {
		return entity.Workload{}, err
	}

	return result, nil
}

func (r *WorkloadRepository) Create(ctx context.Context, value entity.Workload) (entity.Workload, error) {
	if err := coredatabase.TransactionDB(ctx, r.db).Create(&value).Error; err != nil {
		return entity.Workload{}, helpers.WriteError(err)
	}
	return value, nil
}

func (r *WorkloadRepository) ChangeStatus(ctx context.Context, workloadID string, expectedVersion int64, next status.Status) (entity.Workload, error) {
	var value entity.Workload
	result := coredatabase.TransactionDB(ctx, r.db).Raw(queries.ChangeStatus, next, workloadID, expectedVersion, next).Scan(&value)
	if err := helpers.UpdateError(result); err != nil {
		return entity.Workload{}, err
	}
	return value, nil
}
