package impl

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    strategies "grovehub/services/integration/internal/app/services/connection/operations/get/strategies"
)

type getStrategy struct{
    repository connectionrepo.ConnectionRepository
}
var _ strategies.Strategy = (*getStrategy)(nil)
func New(repository connectionrepo.ConnectionRepository) strategies.Strategy {return &getStrategy{repository:repository}}
func (s *getStrategy) Execute(
	ctx context.Context,
	tenantID string,
	id string,
) (connectiondto.Connection, error) {
	if !connectionhelpers.ValidID(tenantID) || !connectionhelpers.ValidID(id) {
		return connectiondto.Connection{}, connectionhelpers.ErrInvalidConnection
	}

	connection, err := s.repository.Get(ctx, tenantID, id)
	if err != nil {
		return connectiondto.Connection{}, err
	}

	return connectionhelpers.ToResponse(connection), nil
}


