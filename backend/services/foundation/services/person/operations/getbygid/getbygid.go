// Package getbygid owns the Person GetByGID use case.
package getbygid

import (
    "context"

    persondto "grovehub/services/foundation/dto/person"
    readstrategies "grovehub/services/foundation/services/person/operations/getbygid/strategies"
)

type GetByGIDOperation struct {
    strategy readstrategies.GetByGIDStrategy
}

func New(strategy readstrategies.GetByGIDStrategy) *GetByGIDOperation { return &GetByGIDOperation{strategy: strategy} }

func (service *GetByGIDOperation) GetByGID(ctx context.Context, request persondto.GetPersonByGIDRequest) (persondto.Person, error) {
    return service.strategy.Execute(ctx, request)
}
