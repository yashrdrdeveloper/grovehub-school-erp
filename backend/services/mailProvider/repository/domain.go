package repository

import (
	"context"
	"errors"

	"grovehub/services/mailProvider/entity"
)

var (
	ErrDomainNotFound      = errors.New("mail provider domain not found")
	ErrDomainAlreadyExists = errors.New("mail provider domain already exists")
)

type DomainRepository interface {
	Create(ctx context.Context, domain entity.Domain) error
	FindAll(ctx context.Context) ([]entity.Domain, error)
	FindById(ctx context.Context, domainId entity.DomainId) (entity.Domain, error)
	FindByName(ctx context.Context, name string) (entity.Domain, error)
	Update(ctx context.Context, domain entity.Domain) error
}
