package impl

import (
    mapper "grovehub/services/integration/internal/app/mappers/connection"
    connectionvalidator "grovehub/services/integration/internal/app/validators/connection"
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    strategies "grovehub/services/integration/internal/app/services/connection/operations/list/strategies"
)

type listStrategy struct{
    mapper mapper.ConnectionMapper
    validator connectionvalidator.ConnectionValidator
    repository connectionrepo.ConnectionRepository
}
var _ strategies.Strategy = (*listStrategy)(nil)
func New(repository connectionrepo.ConnectionRepository, validator connectionvalidator.ConnectionValidator, eventMapper mapper.ConnectionMapper) strategies.Strategy {return &listStrategy{repository:repository,validator:validator,mapper:eventMapper}}
func (s *listStrategy) Execute(
	ctx context.Context,
	tenantID string,
	ownerGID string,
) ([]connectiondto.Connection, error) {
	if err:= s.validator.ValidatePair(tenantID,ownerGID);err!=nil {
		return nil, err
	}

	connections, err := s.repository.List(ctx, tenantID, ownerGID)
	if err != nil {
		return nil, err
	}

	result := make([]connectiondto.Connection, 0, len(connections))
	for _, connection := range connections {
		result = append(result, s.mapper.ToResponse(connection))
	}

	return result, nil
}


