package impl

import (
	"context"
	"strings"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	dto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/repository/role"
)

type CatalogReader struct {
	db *gorm.DB
}

var _ role.CatalogReader = (*CatalogReader)(nil)

func NewCatalogReader(db *gorm.DB) role.CatalogReader {
	return &CatalogReader{db: db}
}

func (reader *CatalogReader) ListActive(ctx context.Context, tenantID string) ([]dto.Role, error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return nil, role.ErrScopeRequired
	}
	if reader.db == nil {
		return nil, transaction.ErrMissingDatabase
	}
	result := make([]dto.Role, 0)
	err := transaction.Database(ctx, reader.db).Raw("SELECT role_id::text AS id, name FROM identity.roles WHERE tenant_id = ? AND status = 'ACTIVE' ORDER BY normalized_name, role_id", tenantID).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
