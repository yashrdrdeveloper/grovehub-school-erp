package impl

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    strategies "grovehub/services/integration/internal/app/services/connection/operations/list/strategies"
)

type listStrategy struct{
    repository connectionrepo.ConnectionRepository
}
var _ strategies.Strategy = (*listStrategy)(nil)
func New(repository connectionrepo.ConnectionRepository) strategies.Strategy {return &listStrategy{repository:repository}}
func (s *listStrategy) Execute(
	ctx context.Context,
	tenantID string,
	ownerGID string,
) ([]connectiondto.Connection, error) {
	if !connectionhelpers.ValidID(tenantID) || !connectionhelpers.ValidID(ownerGID) {
		return nil, connectionhelpers.ErrInvalidConnection
	}

	connections, err := s.repository.List(ctx, tenantID, ownerGID)
	if err != nil {
		return nil, err
	}

	result := make([]connectiondto.Connection, 0, len(connections))
	for _, connection := range connections {
		result = append(result, connectionhelpers.ToResponse(connection))
	}

	return result, nil
}


