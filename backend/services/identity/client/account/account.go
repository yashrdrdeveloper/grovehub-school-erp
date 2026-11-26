package account

import (
	"context"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	dto "grovehub/services/identity/internal/app/dto/account"
)

type (
	Create = dto.Create
	CreateResult = dto.CreateResult
	Account = dto.Account
	Update = dto.Update
	GetRequest = dto.GetRequest
	SearchRequest = dto.SearchRequest
)

// @Client identityAccountClient
type Client interface {
	// @RequestLine GET /protected/accounts/{accountId}
	Get(ctx context.Context, request dto.GetRequest) (dto.Account, error)

	// @RequestLine POST /protected/accounts/search
	Search(ctx context.Context, request dto.SearchRequest) (coresearch.Result[dto.Account], error)
}
