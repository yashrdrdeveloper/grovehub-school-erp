package impl

import (
	"context"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	authorityconstants "grovehub/services/identity/internal/app/repository/authority/constants"
	"grovehub/services/identity/internal/app/repository/registry"
	"grovehub/services/identity/internal/app/repository/registry/helpers"
)

type Repository struct {
	db *gorm.DB
}

var _ registry.Repository = (*Repository)(nil)

func New(db *gorm.DB) registry.Repository {
	return &Repository{db: db}
}

func (repository *Repository) Register(ctx context.Context, batch registry.Batch) (registry.Result, error) {
	var result registry.Result
	err := coredatabase.RunTransaction(ctx, repository.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, repository.db)
		// Registration inserts inactive definitions only; it cannot remove authority.
		if err := db.Exec(authorityconstants.LockGlobalShared).Error; err != nil {
			return err
		}
		for _, definition := range batch.Features {
			created, err := helpers.RegisterFeature(db, definition)
			if err != nil {
				return err
			}
			result.FeaturesCreated += created
		}
		for _, definition := range batch.Permissions {
			created, err := helpers.RegisterPermission(db, definition)
			if err != nil {
				return err
			}
			result.PermissionsCreated += created
		}
		return nil
	})
	if err != nil {
		return registry.Result{}, err
	}
	return result, nil
}
