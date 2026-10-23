package repository

import (
	"context"
	"errors"

	"grovehub/services/identity/entity"
)

var (
	ErrAccountNotFound               = errors.New("account not found")
	ErrAccountAlreadyExists          = errors.New("account already exists")
	ErrBootstrapAlreadyCompleted     = errors.New("identity bootstrap already completed")
	ErrLoginIdentifierNotFound       = errors.New("login identifier not found")
	ErrLoginIdentifierAlreadyExists  = errors.New("login identifier already exists")
	ErrActiveCredentialNotFound      = errors.New("active credential not found")
	ErrActiveCredentialAlreadyExists = errors.New("active credential already exists")
)

type AccountRepository interface {
	HasAccounts(ctx context.Context) (bool, error)
	CreateFirst(
		ctx context.Context,
		account entity.Account,
		identifiers []entity.LoginIdentifier,
		credentials []entity.Credential,
	) error
	Create(
		ctx context.Context,
		account entity.Account,
		identifiers []entity.LoginIdentifier,
		credentials []entity.Credential,
	) error
	FindById(ctx context.Context, accountId entity.AccountId) (entity.Account, error)
	FindByPersonGid(ctx context.Context, personGid entity.PersonGid) (entity.Account, error)
	FindLoginIdentifier(
		ctx context.Context,
		identifierType entity.LoginIdentifierType,
		normalizedValue string,
	) (entity.LoginIdentifier, error)
	FindActiveCredential(
		ctx context.Context,
		accountId entity.AccountId,
		credentialType entity.CredentialType,
	) (entity.Credential, error)
}
