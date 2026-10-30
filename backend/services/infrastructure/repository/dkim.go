package repository

import (
	"context"
	"errors"

	"grovehub/services/infrastructure/entity"
)

var (
	ErrDKIMKeyNotFound      = errors.New("infrastructure DKIM key not found")
	ErrDKIMKeyAlreadyExists = errors.New("infrastructure DKIM key already exists")
)

type DKIMRepository interface {
	Create(ctx context.Context, key entity.DKIMKey) error
	FindByDomain(ctx context.Context, domain string) (entity.DKIMKey, error)
}
