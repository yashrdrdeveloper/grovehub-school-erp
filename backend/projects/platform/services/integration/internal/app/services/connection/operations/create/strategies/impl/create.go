package impl

import (
	"context"
	coreid "grovehub/libs/core/id"
	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	"grovehub/services/integration/internal/app/entity"
	mapper "grovehub/services/integration/internal/app/mappers/connection"
	providercontracts "grovehub/services/integration/internal/app/providers/contracts"
	connectionrepo "grovehub/services/integration/internal/app/repository/connection"
	strategies "grovehub/services/integration/internal/app/services/connection/operations/create/strategies"
	connectionvalidator "grovehub/services/integration/internal/app/validators/connection"
	"time"
)

type createStrategy struct {
	mapper     mapper.ConnectionMapper
	validator  connectionvalidator.ConnectionValidator
	repository connectionrepo.ConnectionRepository
	registry   providercontracts.Registry
}

var _ strategies.Strategy = (*createStrategy)(nil)

func New(repository connectionrepo.ConnectionRepository, registry providercontracts.Registry, validator connectionvalidator.ConnectionValidator, eventMapper mapper.ConnectionMapper) strategies.Strategy {
	return &createStrategy{repository: repository, registry: registry, validator: validator, mapper: eventMapper}
}

func (s *createStrategy) Execute(
	ctx context.Context,
	tenantID string,
	input connectiondto.CreateConnection,
) (connectiondto.Connection, error) {
	if err := s.validator.ValidateCreate(tenantID, input); err != nil {
		return connectiondto.Connection{}, err
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

	return s.mapper.ToResponse(connection), nil
}
