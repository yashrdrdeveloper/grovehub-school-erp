// Package update owns the Person Update use case.
package update

import (
    "context"

    corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    updateflow "grovehub/services/foundation/internal/app/services/person/operations/update/flows"
    updatestrategies "grovehub/services/foundation/internal/app/services/person/operations/update/strategies"
)

type UpdateOperation struct {
    strategy updatestrategies.UpdateStrategy
}

func New(strategy updatestrategies.UpdateStrategy) *UpdateOperation { return &UpdateOperation{strategy: strategy} }

// @ChainOfResponsibility
func (service *UpdateOperation) Update(ctx context.Context, request persondto.UpdatePersonRequest) (persondto.Person, error) {
    result, _, err := corepipeline.RunService(ctx, &updateflow.Update{Request: request},
        service.strategy.Load,
        service.strategy.Apply,
        service.strategy.Persist,
    )
    return result, err
}
