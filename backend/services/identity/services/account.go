package services

import (
	"context"

	"grovehub/services/identity/entity"
)

type CreateAccountInput struct {
	PersonGid string
	Email     string
	Phone     string
}

type UpdateAccountInput struct {
	Email *string
	Phone *string
}

type AccountService interface {
	Create(ctx context.Context, input CreateAccountInput) (entity.Account, error)
	GetById(ctx context.Context, accountId string) (entity.Account, error)
	GetByPersonGid(ctx context.Context, personGid string) (entity.Account, error)
	GetByEmail(ctx context.Context, email string) (entity.Account, error)
	Update(ctx context.Context, accountId string, input UpdateAccountInput) (entity.Account, error)
	Deactivate(ctx context.Context, accountId string) (entity.Account, error)
	PermanentlyDelete(ctx context.Context, accountId string) error
}
