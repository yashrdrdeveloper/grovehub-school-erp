package repository

import (
	"context"
	"errors"

	"grovehub/services/identity/entity"
)

var (
	ErrAccountPersonAlreadyExists = errors.New("identity account already exists for person")
	ErrAccountEmailAlreadyExists  = errors.New("identity account email already exists")
	ErrAccountPhoneAlreadyExists  = errors.New("identity account phone already exists")
)

type AccountRepository interface {
	Create(ctx context.Context, account entity.Account) error
}
