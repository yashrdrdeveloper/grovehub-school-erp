package strategies

import (
    "context"

    coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
    corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    searchflow "grovehub/services/foundation/internal/app/services/person/operations/search/flows"
)

// @Strategy
type SearchStrategy interface {
    Prepare(ctx context.Context, state *searchflow.Search) corepipeline.Outcome[coresearch.Result[persondto.Person]]
    Query(ctx context.Context, state *searchflow.Search) corepipeline.Outcome[coresearch.Result[persondto.Person]]
    Hydrate(ctx context.Context, state *searchflow.Search) corepipeline.Outcome[coresearch.Result[persondto.Person]]
}
