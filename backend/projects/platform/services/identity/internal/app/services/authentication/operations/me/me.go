package me

import (
	"context"
	accountdto "grovehub/services/identity/internal/app/dto/account"
)

// Operation is the use-case contract; concrete dependencies live in impl/.
type Operation interface {
	Execute(ctx context.Context) (accountdto.Account, error)
}
