package flows

import (
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
)

type Search struct {
	Request coresearch.Request
	TenantID string
	Query coresearch.Query[string]
	Accounts coresearch.QueryResult[entity.Account]
	Result coresearch.Result[accountdto.Account]
}
