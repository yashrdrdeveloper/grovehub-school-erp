package account

import (
	"context"

	coresearch "grovehub/libs/core/search"
	dto "grovehub/services/identity/internal/app/dto/account"
)

// Client-facing aliases expose the Identity-owned account DTOs.
type (
	Create = dto.Create
	CreateResult = dto.CreateResult
	Account = dto.Account
	Update = dto.Update
	GetRequest = dto.GetRequest
	ByID = dto.ByID
	SearchRequest = dto.SearchRequest
)

// @Client identityAccountClient
type Client interface {
	// @RequestLine GET /restricted/accounts/{accountId}
	Get(ctx context.Context, request dto.GetRequest) (dto.Account, error)

	// @RequestLine POST /restricted/accounts/search
	Search(ctx context.Context, request dto.SearchRequest) (coresearch.Result[dto.Account], error)
}
