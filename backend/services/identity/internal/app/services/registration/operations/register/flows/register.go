package flows

import (
	foundation "grovehub/services/foundation/client/person"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/security/principal"
)

// Register is scoped to one registration request and is never stored globally.
type Register struct {
	Request accountdto.Create
	IdempotencyKey string
	RequestHash string
	Actor principal.Principal
	Person foundation.Person
	Account entity.Account
	Response accountdto.CreateResult
	Done bool
}
