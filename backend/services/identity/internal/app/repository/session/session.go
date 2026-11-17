package session
import ("context";"errors";"time";"grovehub/services/identity/internal/app/entity")
var ErrInvalidSession=errors.New("invalid or expired session")
type Repository interface{
 Create(context.Context,entity.Session)error
 FindAccess(context.Context,string,time.Time)(entity.Session,error)
 FindRefresh(context.Context,string,time.Time)(entity.Session,error)
 Rotate(context.Context,entity.Session,string)error
 Reauthenticate(context.Context,string,time.Time)error
 List(context.Context,string)([]entity.Session,error)
 Revoke(context.Context,string,string)error
 RevokeAll(context.Context,string)(int64,error)
}
