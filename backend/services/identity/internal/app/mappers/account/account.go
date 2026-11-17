package account

import (
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
)

// @Mapper Mapper
// Explicit manual mapping protects the exported account representation.
// No @Mapping code generation is needed for this domain-aware projection.
type Mapper interface {
	ToResponse(account entity.Account) accountdto.Account
}
