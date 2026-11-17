package passwordtoken
import ("context";"errors";"time";"grovehub/services/identity/internal/app/entity")
var ErrInvalidToken=errors.New("invalid or expired token")
type Repository interface {
 Issue(context.Context,entity.PasswordToken)error
 Consume(context.Context,string,entity.TokenPurpose,time.Time)(string,error)
 Invalidate(context.Context,string,entity.TokenPurpose)error
}
