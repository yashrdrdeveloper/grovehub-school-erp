// Package deactivate owns the Person Deactivate use case.
package deactivate

import (
    "context"

    servicepipeline "grovehub/libs/core/pipeline/service"
    persondto "grovehub/services/foundation/client/person/dto"
    deactivateflow "grovehub/services/foundation/internal/app/services/person/operations/deactivate/flows"
    deactivatestrategies "grovehub/services/foundation/internal/app/services/person/operations/deactivate/strategies"
)

type DeactivateOperation struct { strategy deactivatestrategies.DeactivateStrategy }

func New(strategy deactivatestrategies.DeactivateStrategy) *DeactivateOperation { return &DeactivateOperation{strategy: strategy} }

// @ChainOfResponsibility
func (service *DeactivateOperation) Deactivate(ctx context.Context, request persondto.DeactivatePersonRequest) (persondto.Person, error) {
    result, _, err := servicepipeline.Run(ctx, &deactivateflow.Deactivate{Request: request},
        service.strategy.Load,
        service.strategy.Prepare,
        service.strategy.Persist,
    )
    return result, err
}
