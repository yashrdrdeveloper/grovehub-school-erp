// Package search owns the Person Search use case.
package search

import (
    "context"

    coresearch "grovehub/libs/core/search"
    servicepipeline "grovehub/libs/core/pipeline/service"
    persondto "grovehub/services/foundation/dto/person"
    searchflow "grovehub/services/foundation/services/person/operations/search/flows"
    searchstrategies "grovehub/services/foundation/services/person/operations/search/strategies"
)

type SearchOperation struct { strategy searchstrategies.SearchStrategy }

func New(strategy searchstrategies.SearchStrategy) *SearchOperation { return &SearchOperation{strategy: strategy} }

// @ChainOfResponsibility
func (service *SearchOperation) Search(ctx context.Context, request persondto.SearchPersonsRequest) (coresearch.Result[persondto.Person], error) {
    result, _, err := servicepipeline.Run(ctx, &searchflow.Search{Request: request},
        service.strategy.Prepare,
        service.strategy.Query,
        service.strategy.Hydrate,
    )
    return result, err
}
