package impl

import (
	"context"

	registryrepo "grovehub/services/identity/internal/app/repository/registry"
	"grovehub/services/identity/internal/app/services/registry"
	"grovehub/services/identity/internal/app/services/registry/helpers"
)

// @Service
type Service struct {
	repository registryrepo.Repository
}

var _ registry.Service = (*Service)(nil)

func New(repository registryrepo.Repository) registry.Service {
	return &Service{repository: repository}
}

func (service *Service) Register(ctx context.Context, request registry.Request) (registryrepo.Result, error) {
	batch, err := helpers.PrepareBatch(request)
	if err != nil {
		return registryrepo.Result{}, err
	}
	return service.repository.Register(ctx, batch)
}
