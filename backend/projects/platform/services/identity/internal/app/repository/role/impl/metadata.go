package impl

import (
	"context"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/identity/internal/app/repository/role"
)

type DraftInspector struct {
	db *gorm.DB
}

func NewDraftInspector(db *gorm.DB) role.DraftInspector {
	return &DraftInspector{db: db}
}

func (reader *DraftInspector) Require(ctx context.Context, tenantID string, roleID string, version int64) error {
	if reader.db == nil || !coredatabase.TransactionActive(ctx) {
		return coredatabase.ErrMissingDatabase
	}
	return checkDraft(coredatabase.TransactionDB(ctx, reader.db), role.DraftPermissionUpdate{TenantID: tenantID, RoleID: roleID, Version: version})
}
