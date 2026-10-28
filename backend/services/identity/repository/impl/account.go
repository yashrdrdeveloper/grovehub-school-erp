package impl

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"grovehub/services/identity/entity"
	identityrepository "grovehub/services/identity/repository"
)

const (
	uniqueAccountPersonGidConstraint = "uq_accounts_person_gid"
	uniqueAccountEmailConstraint     = "uq_accounts_email"
	uniqueAccountPhoneConstraint     = "uq_accounts_phone"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) identityrepository.AccountRepository {
	return &AccountRepository{db: db}
}

func (repository *AccountRepository) Create(ctx context.Context, account entity.Account) error {
	return mapAccountWriteError(repository.db.WithContext(ctx).Create(&account).Error)
}

func (repository *AccountRepository) FindById(
	ctx context.Context,
	accountId entity.AccountId,
) (entity.Account, error) {
	var account entity.Account
	if err := repository.db.WithContext(ctx).
		Where(&entity.Account{AccountId: accountId}).
		Take(&account).Error; err != nil {
		return mapAccountReadError(err)
	}
	return account, nil
}

func (repository *AccountRepository) FindByPersonGid(
	ctx context.Context,
	personGid entity.PersonGid,
) (entity.Account, error) {
	var account entity.Account
	if err := repository.db.WithContext(ctx).
		Where(&entity.Account{PersonGid: personGid}).
		Take(&account).Error; err != nil {
		return mapAccountReadError(err)
	}
	return account, nil
}

func (repository *AccountRepository) FindByEmail(
	ctx context.Context,
	email string,
) (entity.Account, error) {
	var account entity.Account
	if err := repository.db.WithContext(ctx).
		Where(&entity.Account{Email: email}).
		Take(&account).Error; err != nil {
		return mapAccountReadError(err)
	}
	return account, nil
}

func (repository *AccountRepository) Update(ctx context.Context, account entity.Account) error {
	result := repository.db.WithContext(ctx).
		Model(&entity.Account{AccountId: account.GetAccountId()}).
		Select("Email", "Phone", "Status", "Version", "UpdatedAt").
		Updates(&account)
	if result.Error != nil {
		return mapAccountWriteError(result.Error)
	}
	if result.RowsAffected == 0 {
		return identityrepository.ErrAccountNotFound
	}
	return nil
}

func (repository *AccountRepository) Delete(
	ctx context.Context,
	accountId entity.AccountId,
) error {
	result := repository.db.WithContext(ctx).
		Delete(&entity.Account{AccountId: accountId})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return identityrepository.ErrAccountNotFound
	}
	return nil
}

func mapAccountReadError(err error) (entity.Account, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Account{}, identityrepository.ErrAccountNotFound
	}
	return entity.Account{}, err
}

func mapAccountWriteError(err error) error {
	if err == nil {
		return nil
	}

	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) || postgresError.Code != "23505" {
		return err
	}

	mappedErrors := map[string]error{
		uniqueAccountPersonGidConstraint: identityrepository.ErrAccountPersonAlreadyExists,
		uniqueAccountEmailConstraint:     identityrepository.ErrAccountEmailAlreadyExists,
		uniqueAccountPhoneConstraint:     identityrepository.ErrAccountPhoneAlreadyExists,
	}
	mappedError, ok := mappedErrors[postgresError.ConstraintName]
	if !ok {
		return err
	}
	return mappedError
}
