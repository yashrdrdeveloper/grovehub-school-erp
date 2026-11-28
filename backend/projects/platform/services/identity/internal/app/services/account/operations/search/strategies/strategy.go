package strategies

import (
	"context"

	corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/services/account/operations/search/flows"
)

// @Strategy
type SearchStrategy interface {
	Prepare(ctx context.Context, state *flows.Search) corepipeline.Outcome[coresearch.Result[accountdto.Account]]
	Query(ctx context.Context, state *flows.Search) corepipeline.Outcome[coresearch.Result[accountdto.Account]]
	Hydrate(ctx context.Context, state *flows.Search) corepipeline.Outcome[coresearch.Result[accountdto.Account]]
}
