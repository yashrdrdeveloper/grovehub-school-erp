// Package search runs tenant-scoped Account Search as an explicit typed chain.
package search

import (
	"context"

	chain "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/service"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/services/account/operations/search/flows"
	strategies "grovehub/services/identity/internal/app/services/account/operations/search/strategies"
)

// Operation is the stable Account Search contract.
type Operation interface {
	Execute(ctx context.Context, request coresearch.Request) (coresearch.Result[accountdto.Account], error)
}

type SearchOperation struct {
	strategy strategies.SearchStrategy
}

var _ Operation = (*SearchOperation)(nil)

func New(strategy strategies.SearchStrategy) Operation {
	return &SearchOperation{strategy: strategy}
}

// @ChainOfResponsibility
func (operation *SearchOperation) Execute(ctx context.Context, request coresearch.Request) (coresearch.Result[accountdto.Account], error) {
	result, _, err := chain.Run(ctx, &flows.Search{Request: request},
		operation.strategy.Prepare,
		operation.strategy.Query,
		operation.strategy.Hydrate,
	)
	return result, err
}
