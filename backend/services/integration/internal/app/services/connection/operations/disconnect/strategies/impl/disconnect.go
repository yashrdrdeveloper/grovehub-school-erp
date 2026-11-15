package impl

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    strategies "grovehub/services/integration/internal/app/services/connection/operations/disconnect/strategies"
)

type disconnectStrategy struct{
    repository connectionrepo.ConnectionRepository
}
var _ strategies.Strategy = (*disconnectStrategy)(nil)
func New(repository connectionrepo.ConnectionRepository) strategies.Strategy {return &disconnectStrategy{repository:repository}}
func (s *disconnectStrategy) Execute(
	ctx context.Context,
	tenantID string,
	id string,
) (connectiondto.Connection, error) {
	if !connectionhelpers.ValidID(tenantID) || !connectionhelpers.ValidID(id) {
		return connectiondto.Connection{}, connectionhelpers.ErrInvalidConnection
	}

	connection, err := s.repository.Disconnect(ctx, tenantID, id)
	if err != nil {
		return connectiondto.Connection{}, err
	}

	return connectionhelpers.ToResponse(connection), nil
}

