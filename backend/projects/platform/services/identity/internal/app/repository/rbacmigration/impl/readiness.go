package impl

import (
	"context"
	"strings"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
)

type Readiness struct {
	db *gorm.DB
}

var _ rbacmigration.Readiness = (*Readiness)(nil)

func NewReadiness(db *gorm.DB) rbacmigration.Readiness {
	return &Readiness{db: db}
}

func (reader *Readiness) Require(ctx context.Context, tenantID string) error {
	if reader.db == nil || strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return rbacmigration.ErrNotReady
	}
	var ready bool
	err := coredatabase.TransactionDB(ctx, reader.db).Raw(
		"SELECT EXISTS (SELECT 1 FROM identity.rbac_migration_receipts b JOIN identity.rbac_migration_receipts s ON s.tenant_id = b.tenant_id AND s.operation = 'BOOTSTRAP' AND s.created_at >= b.created_at WHERE b.tenant_id = ? AND b.operation = 'BACKFILL')", tenantID,
	).Scan(&ready).Error
	if err != nil {
		return err
	}
	if !ready {
		return rbacmigration.ErrNotReady
	}
	return nil
}
