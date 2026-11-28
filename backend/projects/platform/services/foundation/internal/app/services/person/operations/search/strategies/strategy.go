package strategies

import (
	"context"

	corepipeline "grovehub/libs/core/pipeline"
	coresearch "grovehub/libs/core/search"
	persondto "grovehub/services/foundation/internal/app/dto/person"
	searchflow "grovehub/services/foundation/internal/app/services/person/operations/search/flows"
)

// @Strategy
type SearchStrategy interface {
	Prepare(ctx context.Context, state *searchflow.Search) corepipeline.Outcome[coresearch.Result[persondto.Person]]
	Query(ctx context.Context, state *searchflow.Search) corepipeline.Outcome[coresearch.Result[persondto.Person]]
	Hydrate(ctx context.Context, state *searchflow.Search) corepipeline.Outcome[coresearch.Result[persondto.Person]]
}
