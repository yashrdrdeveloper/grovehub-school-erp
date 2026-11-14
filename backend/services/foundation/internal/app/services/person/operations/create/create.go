// Package create owns the Person Create service operation and its chain.
package create

import (
    "context"

    servicepipeline "grovehub/libs/core/pipeline/service"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
    createstrategies "grovehub/services/foundation/internal/app/services/person/operations/create/strategies"
)

// CreateOperation owns the Create-specific strategy dependency; embedding promotes
// Create to the parent PersonService without a forwarding wrapper.
type CreateOperation struct {
    creator createstrategies.CreateStrategy
}

func New(creator createstrategies.CreateStrategy) *CreateOperation { return &CreateOperation{creator: creator} }

// @ChainOfResponsibility
func (service *CreateOperation) Create(ctx context.Context, request persondto.CreatePersonRequest) (persondto.Person, bool, error) {
    flow := personflow.Create{Request: request}
    return servicepipeline.Run(ctx, &flow,
        service.creator.Prepare,
        service.creator.Idempotency,
        service.creator.Contacts,
        service.creator.Persist,
    )
}
