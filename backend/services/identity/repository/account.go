package repository

import (
	"context"
	"errors"

	"grovehub/services/identity/entity"
)

var (
	ErrAccountNotFound            = errors.New("identity account not found")
	ErrAccountPersonAlreadyExists = errors.New("identity account already exists for person")
	ErrAccountEmailAlreadyExists  = errors.New("identity account email already exists")
	ErrAccountPhoneAlreadyExists  = errors.New("identity account phone already exists")
)

type AccountRepository interface {
	Create(ctx context.Context, account entity.Account) error
	FindById(ctx context.Context, accountId entity.AccountId) (entity.Account, error)
	FindByPersonGid(ctx context.Context, personGid entity.PersonGid) (entity.Account, error)
	FindByEmail(ctx context.Context, email string) (entity.Account, error)
	Update(ctx context.Context, account entity.Account) error
	Delete(ctx context.Context, accountId entity.AccountId) error
}
