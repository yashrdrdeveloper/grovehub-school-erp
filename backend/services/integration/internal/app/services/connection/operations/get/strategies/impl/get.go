package impl

import (
    connectionvalidator "grovehub/services/integration/internal/app/validators/connection"
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    strategies "grovehub/services/integration/internal/app/services/connection/operations/get/strategies"
)

type getStrategy struct{
    validator connectionvalidator.ConnectionValidator
    repository connectionrepo.ConnectionRepository
}
var _ strategies.Strategy = (*getStrategy)(nil)
func New(repository connectionrepo.ConnectionRepository, validator connectionvalidator.ConnectionValidator) strategies.Strategy {return &getStrategy{repository:repository,validator:validator}}
func (s *getStrategy) Execute(
	ctx context.Context,
	tenantID string,
	id string,
) (connectiondto.Connection, error) {
	if err:= s.validator.ValidatePair(tenantID,id);err!=nil {
		return connectiondto.Connection{}, err
	}

	connection, err := s.repository.Get(ctx, tenantID, id)
	if err != nil {
		return connectiondto.Connection{}, err
	}

	return connectionhelpers.ToResponse(connection), nil
}


