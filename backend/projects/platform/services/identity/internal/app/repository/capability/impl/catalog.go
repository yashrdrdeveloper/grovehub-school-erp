package impl

import (
	"context"
	"strings"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
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
		return nil, coredatabase.ErrMissingDatabase
	}
	rows := make([]capability.Entry, 0)
	err := coredatabase.TransactionDB(ctx, reader.db).Raw(constants.Available, tenantID).Scan(&rows).Error
	return rows, err
}
