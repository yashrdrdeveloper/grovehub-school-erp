package impl

import (
    mapper "grovehub/services/integration/internal/app/mappers/connection"
    connectionvalidator "grovehub/services/integration/internal/app/validators/connection"
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    strategies "grovehub/services/integration/internal/app/services/connection/operations/disconnect/strategies"
)

type disconnectStrategy struct{
    mapper mapper.ConnectionMapper
    validator connectionvalidator.ConnectionValidator
    repository connectionrepo.ConnectionRepository
}
var _ strategies.Strategy = (*disconnectStrategy)(nil)
func New(repository connectionrepo.ConnectionRepository, validator connectionvalidator.ConnectionValidator, eventMapper mapper.ConnectionMapper) strategies.Strategy {return &disconnectStrategy{repository:repository,validator:validator,mapper:eventMapper}}
func (s *disconnectStrategy) Execute(
	ctx context.Context,
	tenantID string,
	id string,
) (connectiondto.Connection, error) {
	if err:= s.validator.ValidatePair(tenantID,id);err!=nil {
		return connectiondto.Connection{}, err
	}

	connection, err := s.repository.Disconnect(ctx, tenantID, id)
	if err != nil {
		return connectiondto.Connection{}, err
	}

	return s.mapper.ToResponse(connection), nil
}

