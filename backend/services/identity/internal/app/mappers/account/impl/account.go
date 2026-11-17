package impl
import ("grovehub/services/identity/internal/app/dto/account";"grovehub/services/identity/internal/app/entity";contract "grovehub/services/identity/internal/app/mappers/account")
type Mapper struct{}
func New()contract.Mapper{return Mapper{}}
func(Mapper)ToResponse(a entity.Account)account.Account{return account.Account{AccountID:a.AccountID,PersonGID:a.PersonGID,Username:a.Username,Email:a.Email,Status:string(a.Status),Version:a.Version,CreatedAt:a.CreatedAt,UpdatedAt:a.UpdatedAt}}
