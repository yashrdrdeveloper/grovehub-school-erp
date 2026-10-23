package impl

import (
	"context"
	"errors"

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
