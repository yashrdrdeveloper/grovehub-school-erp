package account
import ("context";"errors";"grovehub/libs/core/search";"grovehub/services/identity/internal/app/entity")
var ErrNotFound=errors.New("account not found")
var ErrConflict=errors.New("account conflict")
var ErrVersionConflict=errors.New("account version conflict")
// @Repository entity=Account
type Repository interface {
 Create(context.Context,entity.Account) error
 FindByID(context.Context,string)(entity.Account,error)
 FindByEmail(context.Context,string)(entity.Account,error)
 FindByUsername(context.Context,string)(entity.Account,error)
 FindByPersonGID(context.Context,string)(entity.Account,error)
 FindByIdempotencyKey(context.Context,string)(entity.Account,error)
 Update(context.Context,entity.Account)error
 Search(context.Context,string,search.Query[string])(search.QueryResult[entity.Account],error)
}
