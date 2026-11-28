package impl

import (
	"context"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	authorizationrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/repository/authorization/helpers"
)

type Repository struct {
	db *gorm.DB
}

var _ authorizationrepo.Repository = (*Repository)(nil)

func New(db *gorm.DB) authorizationrepo.Repository {
	return &Repository{db: db}
}

func (repository *Repository) LoadSnapshot(ctx context.Context, tenantID string, accountID string) (authorizationrepo.Snapshot, error) {
	if tenantID == "" || accountID == "" {
		return authorizationrepo.Snapshot{}, authorizationrepo.ErrScopeRequired
	}
	if repository.db == nil {
		return authorizationrepo.Snapshot{}, coredatabase.ErrMissingDatabase
	}
	var rows []helpers.SnapshotRow
	result := coredatabase.TransactionDB(ctx, repository.db).Raw(helpers.SnapshotQuery, tenantID, accountID).Scan(&rows)
	if result.Error != nil {
		return authorizationrepo.Snapshot{}, result.Error
	}
	snapshot, err := helpers.HydrateSnapshot(rows)
	if err != nil {
		return authorizationrepo.Snapshot{}, err
	}
	if snapshot.AccountID != accountID || snapshot.TenantID != tenantID {
		return authorizationrepo.Snapshot{}, authorizationrepo.ErrInvalidSnapshot
	}
	return snapshot, nil
}
