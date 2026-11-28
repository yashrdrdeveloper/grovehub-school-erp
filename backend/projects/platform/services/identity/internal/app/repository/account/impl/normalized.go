package impl

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/account/constants"
)

type NormalizedReader struct {
	db *gorm.DB
}

var _ accountrepo.Reader = (*NormalizedReader)(nil)

func NewNormalizedReader(db *gorm.DB) accountrepo.Reader {
	return &NormalizedReader{db: db}
}

func (reader *NormalizedReader) projection(ctx context.Context) *gorm.DB {
	return coredatabase.TransactionDB(ctx, reader.db).Table("identity.accounts AS a").Select(constants.NormalizedProjection)
}

func (reader *NormalizedReader) FindByID(ctx context.Context, accountID string) (entity.Account, error) {
	id, err := uuid.Parse(accountID)
	if err != nil || id == uuid.Nil {
		return entity.Account{}, accountrepo.ErrNotFound
	}
	if reader.db == nil {
		return entity.Account{}, coredatabase.ErrMissingDatabase
	}
	var account entity.Account
	err = reader.projection(ctx).Where("a.account_id = ?", id.String()).Take(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Account{}, accountrepo.ErrNotFound
	}
	return account, err
}

func (reader *NormalizedReader) Search(ctx context.Context, tenantID string, query coresearch.Query[string]) (coresearch.QueryResult[entity.Account], error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return coresearch.QueryResult[entity.Account]{}, errors.New("account search requires tenant scope")
	}
	if reader.db == nil {
		return coresearch.QueryResult[entity.Account]{}, coredatabase.ErrMissingDatabase
	}
	projection := reader.projection(ctx).Where("a.tenant_id = ?", tenantID)
	db := coredatabase.TransactionDB(ctx, reader.db).Table("(?) AS a", projection)
	return searchAccounts(db, query)
}
