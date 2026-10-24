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
	passwordCredentialPrimaryConstraint   = "pk_password_credentials"
	passwordCredentialAccountFKConstraint = "fk_password_credentials_account"
)

type PasswordCredentialRepository struct {
	pool *pgxpool.Pool
}

func NewPasswordCredentialRepository(pool *pgxpool.Pool) identityrepository.PasswordCredentialRepository {
	return &PasswordCredentialRepository{pool: pool}
}

func (repository *PasswordCredentialRepository) Create(
	ctx context.Context,
	credential entity.PasswordCredential,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`INSERT INTO identity.password_credentials (
			account_id, password_hash, password_changed_at, version, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		string(credential.GetAccountId()),
		credential.GetPasswordHash(),
		credential.GetPasswordChangedAt(),
		credential.GetVersion(),
		credential.GetCreatedAt(),
		credential.GetUpdatedAt(),
	)
	return mapPasswordCredentialWriteError(err)
}

func (repository *PasswordCredentialRepository) FindByAccountId(
	ctx context.Context,
	accountId entity.AccountId,
) (entity.PasswordCredential, error) {
	return scanPasswordCredential(repository.pool.QueryRow(
		ctx,
		passwordCredentialSelect+` WHERE account_id = $1`,
		string(accountId),
	))
}

func (repository *PasswordCredentialRepository) Update(
	ctx context.Context,
	credential entity.PasswordCredential,
) error {
	result, err := repository.pool.Exec(
		ctx,
		`UPDATE identity.password_credentials
		 SET password_hash = $2,
		     password_changed_at = $3,
		     version = $4,
		     updated_at = $5
		 WHERE account_id = $1`,
		string(credential.GetAccountId()),
		credential.GetPasswordHash(),
		credential.GetPasswordChangedAt(),
		credential.GetVersion(),
		credential.GetUpdatedAt(),
	)
	if err != nil {
		return mapPasswordCredentialWriteError(err)
	}
	if result.RowsAffected() == 0 {
		return identityrepository.ErrPasswordCredentialNotFound
	}
	return nil
}

const passwordCredentialSelect = `SELECT
	account_id, password_hash, password_changed_at, version, created_at, updated_at
 FROM identity.password_credentials`

func scanPasswordCredential(row pgx.Row) (entity.PasswordCredential, error) {
	var (
		accountId         string
		passwordHash      string
		passwordChangedAt time.Time
		version           int64
		createdAt         time.Time
		updatedAt         time.Time
	)
	if err := row.Scan(
		&accountId,
		&passwordHash,
		&passwordChangedAt,
		&version,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PasswordCredential{}, identityrepository.ErrPasswordCredentialNotFound
		}
		return entity.PasswordCredential{}, err
	}

	return entity.NewPasswordCredential(
		entity.AccountId(accountId),
		passwordHash,
		passwordChangedAt,
		version,
		createdAt,
		updatedAt,
	), nil
}

func mapPasswordCredentialWriteError(err error) error {
	if err == nil {
		return nil
	}

	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return err
	}

	if postgresError.Code == "23505" && postgresError.ConstraintName == passwordCredentialPrimaryConstraint {
		return identityrepository.ErrPasswordCredentialAlreadyExists
	}
	if postgresError.Code == "23503" && postgresError.ConstraintName == passwordCredentialAccountFKConstraint {
		return identityrepository.ErrPasswordCredentialAccountNotFound
	}
	return err
}
