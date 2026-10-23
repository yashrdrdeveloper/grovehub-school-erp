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
	uniqueAccountPersonGidConstraint       = "uq_accounts_person_gid"
	uniqueLoginIdentifierConstraint        = "uq_login_identifiers_type_normalized_value"
	uniquePrimaryLoginIdentifierConstraint = "uq_login_identifiers_account_primary_type"
	uniqueActiveCredentialConstraint       = "uq_credentials_account_active_type"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) identityrepository.AccountRepository {
	return &AccountRepository{pool: pool}
}

func (repository *AccountRepository) HasAccounts(ctx context.Context) (bool, error) {
	var exists bool
	err := repository.pool.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM identity.accounts)`,
	).Scan(&exists)
	return exists, err
}

func (repository *AccountRepository) CreateFirst(
	ctx context.Context,
	account entity.Account,
	identifiers []entity.LoginIdentifier,
	credentials []entity.Credential,
) error {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `LOCK TABLE identity.accounts IN EXCLUSIVE MODE`); err != nil {
		return err
	}

	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM identity.accounts)`).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return identityrepository.ErrBootstrapAlreadyCompleted
	}

	if err := createAccountAggregate(ctx, tx, account, identifiers, credentials); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (repository *AccountRepository) Create(
	ctx context.Context,
	account entity.Account,
	identifiers []entity.LoginIdentifier,
	credentials []entity.Credential,
) error {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := createAccountAggregate(ctx, tx, account, identifiers, credentials); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (repository *AccountRepository) FindById(
	ctx context.Context,
	accountId entity.AccountId,
) (entity.Account, error) {
	account, err := scanAccount(repository.pool.QueryRow(
		ctx,
		`SELECT account_id, person_gid, status, version, created_at, updated_at
		 FROM identity.accounts
		 WHERE account_id = $1`,
		string(accountId),
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Account{}, identityrepository.ErrAccountNotFound
	}
	return account, err
}

func (repository *AccountRepository) FindByPersonGid(
	ctx context.Context,
	personGid entity.PersonGid,
) (entity.Account, error) {
	account, err := scanAccount(repository.pool.QueryRow(
		ctx,
		`SELECT account_id, person_gid, status, version, created_at, updated_at
		 FROM identity.accounts
		 WHERE person_gid = $1`,
		string(personGid),
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Account{}, identityrepository.ErrAccountNotFound
	}
	return account, err
}

func (repository *AccountRepository) FindLoginIdentifier(
	ctx context.Context,
	identifierType entity.LoginIdentifierType,
	normalizedValue string,
) (entity.LoginIdentifier, error) {
	identifier, err := scanLoginIdentifier(repository.pool.QueryRow(
		ctx,
		`SELECT
			identifier_id, account_id, type, value, normalized_value,
			is_primary, is_verified, status, created_at, updated_at
		 FROM identity.login_identifiers
		 WHERE type = $1 AND normalized_value = $2`,
		string(identifierType),
		normalizedValue,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.LoginIdentifier{}, identityrepository.ErrLoginIdentifierNotFound
	}
	return identifier, err
}

func (repository *AccountRepository) FindActiveCredential(
	ctx context.Context,
	accountId entity.AccountId,
	credentialType entity.CredentialType,
) (entity.Credential, error) {
	credential, err := scanCredential(repository.pool.QueryRow(
		ctx,
		`SELECT
			credential_id, account_id, type, secret_hash, algorithm,
			status, last_used_at, created_at, updated_at
		 FROM identity.credentials
		 WHERE account_id = $1 AND type = $2 AND status = 'ACTIVE'`,
		string(accountId),
		string(credentialType),
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Credential{}, identityrepository.ErrActiveCredentialNotFound
	}
	return credential, err
}

func createAccountAggregate(
	ctx context.Context,
	tx pgx.Tx,
	account entity.Account,
	identifiers []entity.LoginIdentifier,
	credentials []entity.Credential,
) error {
	if err := insertAccount(ctx, tx, account); err != nil {
		return err
	}
	for _, identifier := range identifiers {
		if err := insertLoginIdentifier(ctx, tx, identifier); err != nil {
			return err
		}
	}
	for _, credential := range credentials {
		if err := insertCredential(ctx, tx, credential); err != nil {
			return err
		}
	}
	return nil
}

func insertAccount(ctx context.Context, tx pgx.Tx, account entity.Account) error {
	_, err := tx.Exec(
		ctx,
		`INSERT INTO identity.accounts (
			account_id, person_gid, status, version, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		string(account.GetAccountId()),
		string(account.GetPersonGid()),
		string(account.GetStatus()),
		account.GetVersion(),
		account.GetCreatedAt(),
		account.GetUpdatedAt(),
	)
	return mapAccountWriteError(err)
}

func insertLoginIdentifier(
	ctx context.Context,
	tx pgx.Tx,
	identifier entity.LoginIdentifier,
) error {
	_, err := tx.Exec(
		ctx,
		`INSERT INTO identity.login_identifiers (
			identifier_id, account_id, type, value, normalized_value,
			is_primary, is_verified, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(identifier.GetIdentifierId()),
		string(identifier.GetAccountId()),
		string(identifier.GetType()),
		identifier.GetValue(),
		identifier.GetNormalizedValue(),
		identifier.GetIsPrimary(),
		identifier.GetIsVerified(),
		string(identifier.GetStatus()),
		identifier.GetCreatedAt(),
		identifier.GetUpdatedAt(),
	)
	return mapAccountWriteError(err)
}

func insertCredential(ctx context.Context, tx pgx.Tx, credential entity.Credential) error {
	_, err := tx.Exec(
		ctx,
		`INSERT INTO identity.credentials (
			credential_id, account_id, type, secret_hash, algorithm,
			status, last_used_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		string(credential.GetCredentialId()),
		string(credential.GetAccountId()),
		string(credential.GetType()),
		credential.GetSecretHash(),
		string(credential.GetAlgorithm()),
		string(credential.GetStatus()),
		credential.GetLastUsedAt(),
		credential.GetCreatedAt(),
		credential.GetUpdatedAt(),
	)
	return mapAccountWriteError(err)
}

func scanAccount(row pgx.Row) (entity.Account, error) {
	var (
		accountId entity.AccountId
		personGid entity.PersonGid
		status    entity.AccountStatus
		version   int64
		createdAt time.Time
		updatedAt time.Time
	)
	if err := row.Scan(
		&accountId,
		&personGid,
		&status,
		&version,
		&createdAt,
		&updatedAt,
	); err != nil {
		return entity.Account{}, err
	}
	return entity.NewAccount(accountId, personGid, status, version, createdAt, updatedAt), nil
}

func scanLoginIdentifier(row pgx.Row) (entity.LoginIdentifier, error) {
	var (
		identifierId    entity.LoginIdentifierId
		accountId       entity.AccountId
		identifierType  entity.LoginIdentifierType
		value           string
		normalizedValue string
		isPrimary       bool
		isVerified      bool
		status          entity.LoginIdentifierStatus
		createdAt       time.Time
		updatedAt       time.Time
	)
	if err := row.Scan(
		&identifierId,
		&accountId,
		&identifierType,
		&value,
		&normalizedValue,
		&isPrimary,
		&isVerified,
		&status,
		&createdAt,
		&updatedAt,
	); err != nil {
		return entity.LoginIdentifier{}, err
	}
	return entity.NewLoginIdentifier(
		identifierId,
		accountId,
		identifierType,
		value,
		normalizedValue,
		isPrimary,
		isVerified,
		status,
		createdAt,
		updatedAt,
	), nil
}

func scanCredential(row pgx.Row) (entity.Credential, error) {
	var (
		credentialId   entity.CredentialId
		accountId      entity.AccountId
		credentialType entity.CredentialType
		secretHash     string
		algorithm      entity.CredentialAlgorithm
		status         entity.CredentialStatus
		lastUsedAt     *time.Time
		createdAt      time.Time
		updatedAt      time.Time
	)
	if err := row.Scan(
		&credentialId,
		&accountId,
		&credentialType,
		&secretHash,
		&algorithm,
		&status,
		&lastUsedAt,
		&createdAt,
		&updatedAt,
	); err != nil {
		return entity.Credential{}, err
	}
	return entity.NewCredential(
		credentialId,
		accountId,
		credentialType,
		secretHash,
		algorithm,
		status,
		lastUsedAt,
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
		uniqueAccountPersonGidConstraint:       identityrepository.ErrAccountAlreadyExists,
		uniqueLoginIdentifierConstraint:        identityrepository.ErrLoginIdentifierAlreadyExists,
		uniquePrimaryLoginIdentifierConstraint: identityrepository.ErrLoginIdentifierAlreadyExists,
		uniqueActiveCredentialConstraint:       identityrepository.ErrActiveCredentialAlreadyExists,
	}
	mappedError, ok := mappedErrors[postgresError.ConstraintName]
	if !ok {
		return err
	}
	return mappedError
}
