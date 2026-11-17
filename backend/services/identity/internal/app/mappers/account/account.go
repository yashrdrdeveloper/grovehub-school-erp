package account
import ("grovehub/services/identity/internal/app/dto/account";"grovehub/services/identity/internal/app/entity")
type Mapper interface{ToResponse(entity.Account)account.Account}
