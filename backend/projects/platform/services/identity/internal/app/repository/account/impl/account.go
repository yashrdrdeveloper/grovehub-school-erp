package impl

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	searchconfig "grovehub/services/identity/internal/app/repository/account/helpers/search"
)

type Repository struct {
	db *gorm.DB
}

var _ accountrepo.Repository = (*Repository)(nil)

func New(db *gorm.DB) accountrepo.Repository {
	return &Repository{db: db}
}

func (repository *Repository) dbFor(ctx context.Context) *gorm.DB {
	return coredatabase.TransactionDB(ctx, repository.db)
}

func (repository *Repository) Create(ctx context.Context, account entity.Account) error {
	return mapWriteError(repository.dbFor(ctx).Create(&account).Error)
}

func (repository *Repository) find(ctx context.Context, field string, value any) (entity.Account, error) {
	var account entity.Account
	err := repository.dbFor(ctx).Where(field+" = ?", value).Take(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Account{}, accountrepo.ErrNotFound
	}
	return account, err
}

func (repository *Repository) FindByID(ctx context.Context, accountID string) (entity.Account, error) {
	return repository.find(ctx, "account_id", accountID)
}

func (repository *Repository) FindByEmail(ctx context.Context, email string) (entity.Account, error) {
	return repository.find(ctx, "email", strings.ToLower(email))
}

func (repository *Repository) FindByUsername(ctx context.Context, username string) (entity.Account, error) {
	return repository.find(ctx, "lower(username)", strings.ToLower(username))
}

func (repository *Repository) FindByPersonGID(ctx context.Context, gid string) (entity.Account, error) {
	return repository.find(ctx, "person_gid", gid)
}

func (repository *Repository) FindByIdempotencyKey(ctx context.Context, key string) (entity.Account, error) {
	return repository.find(ctx, "registration_idempotency_key", key)
}

func (repository *Repository) Update(ctx context.Context, account entity.Account) error {
	secondary := string(account.SecondaryRoleIDs)
	if secondary == "" || secondary == "null" {
		secondary = "[]"
	}
	result := repository.dbFor(ctx).
		Model(&entity.Account{}).
		Where("account_id = ? AND version = ?", account.AccountID, account.Version-1).
		Updates(map[string]any{
			"username":           account.Username,
			"primary_role_id":    account.PrimaryRoleID,
			"secondary_role_ids": gorm.Expr("?::jsonb", secondary),
			"email":              account.Email,
			"status":             account.Status,
			"version":            account.Version,
			"search_details":     account.SearchDetails,
			"updated_at":         account.UpdatedAt,
		})
	if result.Error != nil {
		return mapWriteError(result.Error)
	}
	if result.RowsAffected != 1 {
		return accountrepo.ErrVersionConflict
	}
	return nil
}

// Search uses a trusted SearchDetails configuration. The tenant predicate is
// required here because a search generator without scope would leak accounts.
func (repository *Repository) Search(ctx context.Context, tenantID string, query coresearch.Query[string]) (coresearch.QueryResult[entity.Account], error) {
	if tenantID == "" {
		return coresearch.QueryResult[entity.Account]{}, errors.New("account search requires tenant scope")
	}

	db := repository.dbFor(ctx).
		Table("identity.accounts AS a").
		Where("a.tenant_id = ?", tenantID)

	return searchAccounts(db, query)
}

func searchAccounts(db *gorm.DB, query coresearch.Query[string]) (coresearch.QueryResult[entity.Account], error) {
	if query.PageNumber > 0 {
		return coresearch.Execute[string, entity.Account](db, query, searchconfig.AccountConfig)
	}
	// Legacy cursor-only requests still receive the filtered record count.
	countQuery := query
	countQuery.Cursor = nil
	countQuery.Offset = 0
	countDB, err := coresearch.Apply(db.Session(&gorm.Session{}), countQuery, searchconfig.AccountConfig)
	if err != nil {
		return coresearch.QueryResult[entity.Account]{}, err
	}
	var total int64
	if err := countDB.Model(&entity.Account{}).Select("COUNT(*)").Limit(-1).Offset(-1).Count(&total).Error; err != nil {
		return coresearch.QueryResult[entity.Account]{}, err
	}
	result, err := coresearch.Execute[string, entity.Account](db, query, searchconfig.AccountConfig)
	if err != nil {
		return coresearch.QueryResult[entity.Account]{}, err
	}
	result.TotalRecords = &total
	return result, nil
}

func mapWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgres *pgconn.PgError
	if errors.As(err, &postgres) && postgres.Code == "23505" {
		return accountrepo.ErrConflict
	}
	return err
}
