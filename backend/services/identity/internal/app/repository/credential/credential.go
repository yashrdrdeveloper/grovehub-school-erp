package credential
import ("context";"errors";"grovehub/services/identity/internal/app/entity")
var ErrNotFound=errors.New("credential not found")
type Repository interface {
 Find(context.Context,string)(entity.PasswordCredential,error)
 Create(context.Context,entity.PasswordCredential)error
 Update(context.Context,entity.PasswordCredential)error
}
