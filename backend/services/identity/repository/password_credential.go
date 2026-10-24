package repository

import (
	"context"
	"errors"

	"grovehub/services/identity/entity"
)

var (
	ErrPasswordCredentialNotFound       = errors.New("password credential not found")
	ErrPasswordCredentialAlreadyExists  = errors.New("password credential already exists")
	ErrPasswordCredentialAccountNotFound = errors.New("identity account not found for password credential")
)

type PasswordCredentialRepository interface {
	Create(ctx context.Context, credential entity.PasswordCredential) error
	FindByAccountId(ctx context.Context, accountId entity.AccountId) (entity.PasswordCredential, error)
	Update(ctx context.Context, credential entity.PasswordCredential) error
}
