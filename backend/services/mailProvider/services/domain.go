package services

import (
	"context"

	"grovehub/services/mailProvider/entity"
)

type RegisterDomainInput struct {
	Name string
}

type DomainService interface {
	Register(ctx context.Context, input RegisterDomainInput) (entity.Domain, error)
	GetById(ctx context.Context, domainId string) (entity.Domain, error)
	Verify(ctx context.Context, domainId string) (entity.Domain, error)
	Deactivate(ctx context.Context, domainId string) (entity.Domain, error)
}
