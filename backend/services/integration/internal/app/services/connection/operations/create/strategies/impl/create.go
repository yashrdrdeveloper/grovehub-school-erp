package impl

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "time"
    coreid "grovehub/libs/core/id"
    "grovehub/services/integration/internal/app/entity"
    "grovehub/services/integration/internal/app/providers"
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    strategies "grovehub/services/integration/internal/app/services/connection/operations/create/strategies"
)

type createStrategy struct{
    repository connectionrepo.ConnectionRepository
    registry providers.Registry
}
var _ strategies.Strategy = (*createStrategy)(nil)
func New(repository connectionrepo.ConnectionRepository, registry providers.Registry) strategies.Strategy {return &createStrategy{repository:repository,registry:registry}}
func (s *createStrategy) Execute(
	ctx context.Context,
	tenantID string,
	input connectiondto.CreateConnection,
) (connectiondto.Connection, error) {
	if !connectionhelpers.ValidID(tenantID) || !connectionhelpers.ValidID(input.OwnerGID) {
		return connectiondto.Connection{}, connectionhelpers.ErrInvalidConnection
	}

	provider, err := s.registry.Find(input.Provider)
	if err != nil {
		return connectiondto.Connection{}, err
	}

	id, err := coreid.NewUUIDv7()
	if err != nil {
		return connectiondto.Connection{}, err
	}

	now := time.Now().UTC()
	connection := entity.Connection{
		ID:        id,
		TenantID:  tenantID,
		OwnerGID:  input.OwnerGID,
		Provider:  provider.Descriptor().ID,
		Status:    entity.StatusPendingSetup,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(ctx, connection); err != nil {
		return connectiondto.Connection{}, err
	}

	return connectionhelpers.ToResponse(connection), nil
}


