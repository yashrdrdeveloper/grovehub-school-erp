package impl

import (
	"context"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/repository/role"
)

type DraftInspector struct {
	db *gorm.DB
}

func NewDraftInspector(db *gorm.DB) role.DraftInspector {
	return &DraftInspector{db: db}
}

func (reader *DraftInspector) Require(ctx context.Context, tenantID string, roleID string, version int64) error {
	if reader.db == nil || !transaction.Active(ctx) {
		return transaction.ErrMissingDatabase
	}
	return checkDraft(transaction.Database(ctx, reader.db), role.DraftPermissionUpdate{TenantID: tenantID, RoleID: roleID, Version: version})
}
