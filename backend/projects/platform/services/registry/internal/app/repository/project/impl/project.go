package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/constants"
	"grovehub/services/registry/internal/app/repository/project"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) project.Repository {
	return &repository{db: db}
}

func (r *repository) Find(ctx context.Context, projectID string) (entity.Project, error) {
	var result entity.Project
	err := coredatabase.TransactionDB(ctx, r.db).Where("project_id = ?", projectID).Take(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Project{}, constants.ErrNotFound
	}
	if err != nil {
		return entity.Project{}, err
	}

	return result, nil
}
