package impl

import (
	"context"
	"strings"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/authority/constants"
)

type AdministratorReader struct {
	db *gorm.DB
}

var _ authority.AdministratorReader = (*AdministratorReader)(nil)

func NewAdministratorReader(db *gorm.DB) authority.AdministratorReader {
	return &AdministratorReader{db: db}
}

func (reader *AdministratorReader) Count(ctx context.Context, tenantID string) (int64, error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return 0, authority.ErrScopeRequired
	}
	if reader.db == nil {
		return 0, coredatabase.ErrMissingDatabase
	}
	required := permission.AdministratorRecoveryCapabilities()
	codes := make([]string, 0, len(required))
	for _, code := range required {
		definition, known := permission.Find(code)
		if !known || definition.Feature != feature.RoleManagement {
			return 0, authority.ErrBootstrapRequired
		}
		codes = append(codes, string(code))
	}
	var count int64
	err := coredatabase.TransactionDB(ctx, reader.db).Raw(constants.CountAdministrators, tenantID, feature.RoleManagement, codes, len(codes)).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
