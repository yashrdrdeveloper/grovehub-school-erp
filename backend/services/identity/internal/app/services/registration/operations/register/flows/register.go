package flows
import(
 dto "grovehub/services/identity/internal/app/dto/account"
 "grovehub/services/identity/internal/app/entity"
 "grovehub/services/identity/internal/app/security/principal"
 foundation "grovehub/services/foundation/client/person"
)
type Register struct{
 Request dto.Create
 IdempotencyKey string
 RequestHash string
 Actor principal.Principal
 Person foundation.Person
 Account entity.Account
 Response dto.CreateResult
 Done bool
}
