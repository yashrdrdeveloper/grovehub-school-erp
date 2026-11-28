package impl

import (
	"context"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/authority/constants"
)

type GlobalGuard struct {
	db *gorm.DB
	reader authority.AdministratorReader
}

var _ authority.GlobalGuard = (*GlobalGuard)(nil)

func NewGlobalGuard(db *gorm.DB) authority.GlobalGuard {
	return &GlobalGuard{db: db, reader: NewAdministratorReader(db)}
}

func (guard *GlobalGuard) Run(ctx context.Context, operation func(context.Context) error) error {
	if guard.db == nil || guard.reader == nil || operation == nil {
		return authority.ErrScopeRequired
	}
	// Never upgrade a shared global lock held by a tenant transaction.
	if coredatabase.TransactionActive(ctx) {
		return authority.ErrLockOrder
	}
	return coredatabase.RunTransaction(ctx, guard.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, guard.db)
		if err := db.Exec(constants.LockGlobalExclusive).Error; err != nil {
			return err
		}
		protected, err := guard.protectedTenants(tx)
		if err != nil {
			return err
		}
		if err := operation(tx); err != nil {
			return err
		}
		for _, tenantID := range protected {
			count, err := guard.reader.Count(tx, tenantID)
			if err != nil {
				return err
			}
			if count < 1 {
				return authority.ErrLastAdministrator
			}
		}
		return nil
	})
}

func (guard *GlobalGuard) protectedTenants(ctx context.Context) ([]string, error) {
	db := coredatabase.TransactionDB(ctx, guard.db)
	var tenants []string
	if err := db.Raw(constants.Tenants).Scan(&tenants).Error; err != nil {
		return nil, err
	}
	protected := make([]string, 0, len(tenants))
	for _, tenantID := range tenants {
		count, err := guard.reader.Count(ctx, tenantID)
		if err != nil {
			return nil, err
		}
		var bootstrapped bool
		if err := db.Raw(constants.Bootstrapped, tenantID).Scan(&bootstrapped).Error; err != nil {
			return nil, err
		}
		if bootstrapped && count < 1 {
			return nil, authority.ErrBootstrapRequired
		}
		if count > 0 {
			protected = append(protected, tenantID)
		}
	}
	return protected, nil
}
