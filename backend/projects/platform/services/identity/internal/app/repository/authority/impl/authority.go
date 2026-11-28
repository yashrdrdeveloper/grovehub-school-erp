package impl

import (
	"context"
	"strings"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/authority/constants"
)

const lockTenant = `SELECT pg_advisory_xact_lock(hashtextextended('identity:authority:' || CAST(? AS text), 0))`

type Locker struct {
	db *gorm.DB
}

var _ authority.Locker = (*Locker)(nil)

func New(db *gorm.DB) authority.Locker {
	return &Locker{db: db}
}

func (locker *Locker) Run(ctx context.Context, tenantID string, operation func(context.Context) error) error {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID || operation == nil {
		return authority.ErrScopeRequired
	}
	return coredatabase.RunTransaction(ctx, locker.db, func(tx context.Context) error {
		if err := coredatabase.TransactionDB(tx, locker.db).Exec(constants.LockGlobalShared).Error; err != nil {
			return err
		}
		if err := coredatabase.TransactionDB(tx, locker.db).Exec(lockTenant, tenantID).Error; err != nil {
			return err
		}
		return operation(tx)
	})
}
