package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/constants"
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
