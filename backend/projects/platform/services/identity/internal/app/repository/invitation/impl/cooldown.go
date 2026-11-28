package impl

import (
	"context"
	"time"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/identity/internal/app/repository/invitation"
)

type CooldownReader struct {
	db *gorm.DB
}

func NewCooldownReader(db *gorm.DB) invitation.CooldownReader {
	return &CooldownReader{db: db}
}

func (reader *CooldownReader) RecentSetup(ctx context.Context, accountID string, since time.Time) (bool, error) {
	if reader.db == nil || !coredatabase.TransactionActive(ctx) {
		return false, coredatabase.ErrMissingDatabase
	}
	var row struct { Recent bool }
	err := coredatabase.TransactionDB(ctx, reader.db).Raw("SELECT EXISTS (SELECT 1 FROM identity.invitations WHERE account_id = ? AND purpose = 'ACCOUNT_SETUP' AND created_at > ?) AS recent", accountID, since).Scan(&row).Error
	return row.Recent, err
}
