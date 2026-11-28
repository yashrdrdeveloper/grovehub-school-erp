package strategies

import (
    "context"

    coresearch "grovehub/libs/core/search"
    corepipeline "grovehub/libs/core/pipeline"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    searchflow "grovehub/services/foundation/internal/app/services/person/operations/search/flows"
)

// @Strategy
type SearchStrategy interface {
    Prepare(ctx context.Context, state *searchflow.Search) corepipeline.Outcome[coresearch.Result[persondto.Person]]
    Query(ctx context.Context, state *searchflow.Search) corepipeline.Outcome[coresearch.Result[persondto.Person]]
    Hydrate(ctx context.Context, state *searchflow.Search) corepipeline.Outcome[coresearch.Result[persondto.Person]]
}
