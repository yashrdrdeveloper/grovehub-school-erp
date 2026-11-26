package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"gorm.io/gorm"

	"grovehub/services/identity/internal/app/repository/rbacmigration/constants"
	"grovehub/services/identity/internal/app/services/rbacmigration"
)

type Reader struct {
	db *gorm.DB
}

var _ rbacmigration.Reader = (*Reader)(nil)

func New(db *gorm.DB) rbacmigration.Reader {
	return &Reader{db: db}
}

type accountRow struct {
	AccountID string
	Version int64
	Primary *string
	SecondaryJSON string
}

func (reader *Reader) Inventory(ctx context.Context, tenantID string) (rbacmigration.Inventory, error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return rbacmigration.Inventory{}, rbacmigration.ErrScopeRequired
	}
	if reader.db == nil {
		return rbacmigration.Inventory{}, rbacmigration.ErrInvalidInventory
	}
	inventory := rbacmigration.Inventory{TenantID: tenantID}
	err := reader.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rows []accountRow
		if err := tx.Raw(constants.Accounts, tenantID).Scan(&rows).Error; err != nil {
			return err
		}
		for _, row := range rows {
			account := rbacmigration.LegacyAccount{AccountID: row.AccountID, Version: row.Version, Primary: row.Primary}
			if err := json.Unmarshal([]byte(row.SecondaryJSON), &account.Secondary); err != nil {
				return rbacmigration.ErrInvalidInventory
			}
			inventory.Accounts = append(inventory.Accounts, account)
		}
		if err := tx.Raw(constants.Roles, tenantID).Scan(&inventory.Roles).Error; err != nil {
			return err
		}
		return tx.Raw(constants.AssignmentCount, tenantID).Scan(&inventory.NormalizedAssignments).Error
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return rbacmigration.Inventory{}, err
	}
	return inventory, nil
}
