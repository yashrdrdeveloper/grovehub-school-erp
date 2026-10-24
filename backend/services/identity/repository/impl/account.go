package impl

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"grovehub/services/identity/entity"
	identityrepository "grovehub/services/identity/repository"
)

const (
	uniqueAccountPersonGidConstraint = "uq_accounts_person_gid"
	uniqueAccountEmailConstraint     = "uq_accounts_email"
	uniqueAccountPhoneConstraint     = "uq_accounts_phone"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) identityrepository.AccountRepository {
	return &AccountRepository{pool: pool}
}

func (repository *AccountRepository) Create(ctx context.Context, account entity.Account) error {
	_, err := repository.pool.Exec(
		ctx,
		`INSERT INTO identity.accounts (
			account_id, person_gid, email, phone, status, version, created_at, updated_at
		) VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7, $8)`,
		string(account.GetAccountId()),
		string(account.GetPersonGid()),
		account.GetEmail(),
		account.GetPhone(),
		string(account.GetStatus()),
		account.GetVersion(),
		account.GetCreatedAt(),
		account.GetUpdatedAt(),
	)
	return mapAccountWriteError(err)
}

func (repository *AccountRepository) FindById(
	ctx context.Context,
	accountId entity.AccountId,
) (entity.Account, error) {
	return scanAccount(repository.pool.QueryRow(
		ctx,
		accountSelect+` WHERE account_id = $1`,
		string(accountId),
	))
}

func (repository *AccountRepository) FindByPersonGid(
	ctx context.Context,
	personGid entity.PersonGid,
) (entity.Account, error) {
	return scanAccount(repository.pool.QueryRow(
		ctx,
		accountSelect+` WHERE person_gid = $1`,
		string(personGid),
	))
}

func (repository *AccountRepository) FindByEmail(
	ctx context.Context,
	email string,
) (entity.Account, error) {
	return scanAccount(repository.pool.QueryRow(
		ctx,
		accountSelect+` WHERE email = $1`,
		email,
	))
}

func (repository *AccountRepository) Update(ctx context.Context, account entity.Account) error {
	result, err := repository.pool.Exec(
		ctx,
		`UPDATE identity.accounts
		 SET email = $2,
		     phone = NULLIF($3, ''),
		     status = $4,
		     version = $5,
		     updated_at = $6
		 WHERE account_id = $1`,
		string(account.GetAccountId()),
		account.GetEmail(),
		account.GetPhone(),
		string(account.GetStatus()),
		account.GetVersion(),
		account.GetUpdatedAt(),
	)
	if err != nil {
		return mapAccountWriteError(err)
	}
	if result.RowsAffected() == 0 {
		return identityrepository.ErrAccountNotFound
	}
	return nil
}

func (repository *AccountRepository) Delete(
	ctx context.Context,
	accountId entity.AccountId,
) error {
	result, err := repository.pool.Exec(
		ctx,
		`DELETE FROM identity.accounts WHERE account_id = $1`,
		string(accountId),
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return identityrepository.ErrAccountNotFound
	}
	return nil
}

const accountSelect = `SELECT
	account_id, person_gid, email, phone, status, version, created_at, updated_at
 FROM identity.accounts`

func scanAccount(row pgx.Row) (entity.Account, error) {
	var (
		accountId string
		personGid string
		email     string
		phone     *string
		status    string
		version   int64
		createdAt time.Time
		updatedAt time.Time
	)
	if err := row.Scan(
		&accountId,
		&personGid,
		&email,
		&phone,
		&status,
		&version,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Account{}, identityrepository.ErrAccountNotFound
		}
		return entity.Account{}, err
	}

	phoneValue := ""
	if phone != nil {
		phoneValue = *phone
	}
	return entity.NewAccount(
		entity.AccountId(accountId),
		entity.PersonGid(personGid),
		email,
		phoneValue,
		entity.AccountStatus(status),
		version,
		createdAt,
		updatedAt,
	), nil
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
