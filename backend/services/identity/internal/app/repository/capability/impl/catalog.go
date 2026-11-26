package impl

import (
	"context"
	"strings"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/capability"
	"grovehub/services/identity/internal/app/repository/capability/constants"
)

type Reader struct {
	db *gorm.DB
}

func New(db *gorm.DB) capability.Reader {
	return &Reader{db: db}
}

func (reader *Reader) Available(ctx context.Context, tenantID string) ([]capability.Entry, error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return nil, authority.ErrScopeRequired
	}
	if reader.db == nil {
		return nil, transaction.ErrMissingDatabase
	}
	rows := make([]capability.Entry, 0)
	err := transaction.Database(ctx, reader.db).Raw(constants.Available, tenantID).Scan(&rows).Error
	return rows, err
}
