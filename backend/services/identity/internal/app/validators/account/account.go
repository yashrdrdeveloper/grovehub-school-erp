package account

import (
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
)

// Validator exposes pure normalization and trusted SearchDetails query rules.
// Resource authorization remains in its owning service operation.
type Validator interface {
	NormalizeCreate(request accountdto.Create) (accountdto.Create, error)
	NormalizeUsername(username string) (string, error)
	Search(request coresearch.Request) (coresearch.Query[string], error)
}
