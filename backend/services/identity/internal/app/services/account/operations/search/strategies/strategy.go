package strategies

import (
	"context"

	chain "grovehub/libs/core/pipeline/service"
	coresearch "grovehub/libs/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/services/account/operations/search/flows"
)

// @Strategy
type SearchStrategy interface {
	Prepare(ctx context.Context, state *flows.Search) chain.Outcome[coresearch.Result[accountdto.Account]]
	Query(ctx context.Context, state *flows.Search) chain.Outcome[coresearch.Result[accountdto.Account]]
	Hydrate(ctx context.Context, state *flows.Search) chain.Outcome[coresearch.Result[accountdto.Account]]
}
