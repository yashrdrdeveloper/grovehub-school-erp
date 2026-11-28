package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/constants"
	"grovehub/services/registry/internal/app/repository/helpers"
	"grovehub/services/registry/internal/app/repository/project"
	queries "grovehub/services/registry/internal/app/repository/project/constants"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) project.ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Find(ctx context.Context, projectID string) (entity.Project, error) {
	var result entity.Project
	err := coredatabase.TransactionDB(ctx, r.db).Where(queries.FindByID, projectID).Take(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Project{}, constants.ErrNotFound
	}
	if err != nil {
		return entity.Project{}, err
	}

	return result, nil
}

func (r *ProjectRepository) Create(ctx context.Context, value entity.Project) (entity.Project, error) {
	if err := coredatabase.TransactionDB(ctx, r.db).Create(&value).Error; err != nil {
		return entity.Project{}, helpers.WriteError(err)
	}
	return value, nil
}

func (r *ProjectRepository) ChangeStatus(ctx context.Context, projectID string, expectedVersion int64, next status.Status) (entity.Project, error) {
	var value entity.Project
	result := coredatabase.TransactionDB(ctx, r.db).Raw(queries.ChangeStatus, next, projectID, expectedVersion, next).Scan(&value)
	if err := helpers.UpdateError(result); err != nil {
		return entity.Project{}, err
	}
	return value, nil
}
