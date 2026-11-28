package impl

import (
	"context"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/libs/platform/scope"
	client "grovehub/services/registry/client/authorization"
	"grovehub/services/registry/internal/app/repository/authorization"
	queries "grovehub/services/registry/internal/app/repository/authorization/constants"
	"grovehub/services/registry/internal/app/repository/constants"
)

type AuthorizationRepository struct {
	db *gorm.DB
}

func NewAuthorizationRepository(db *gorm.DB) authorization.AuthorizationRepository {
	return &AuthorizationRepository{db: db}
}

func (r *AuthorizationRepository) Resolve(ctx context.Context, keyID string, target scope.Target, capability string, allowInactiveTarget bool) (client.Credential, error) {
	var value client.Credential
	result := coredatabase.TransactionDB(ctx, r.db).Raw(queries.Resolve, keyID, target.Level, capability, target.ProjectID, target.TenantID, allowInactiveTarget).Scan(&value)
	if result.Error != nil {
		return client.Credential{}, result.Error
	}
	if result.RowsAffected != 1 {
		return client.Credential{}, constants.ErrNotFound
	}
	return value, nil
}

func (r *AuthorizationRepository) TenantState(ctx context.Context, projectID, tenantID string) (client.TenantState, error) {
	var value client.TenantState
	result := coredatabase.TransactionDB(ctx, r.db).Raw(queries.TenantState, projectID, tenantID).Scan(&value)
	if result.Error != nil {
		return client.TenantState{}, result.Error
	}
	if result.RowsAffected != 1 || value.SessionGeneration < 1 {
		return client.TenantState{}, constants.ErrNotFound
	}
	return value, nil
}
