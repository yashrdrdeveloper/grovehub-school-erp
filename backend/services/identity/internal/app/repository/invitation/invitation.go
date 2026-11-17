package invitation
import ("context";"time";"grovehub/services/identity/internal/app/entity")
type Repository interface {
 Queue(context.Context,entity.InvitationMail)error
 Claim(context.Context,time.Time)(entity.InvitationMail,bool,error)
 Delivered(context.Context,string)error
 Retry(context.Context,string,time.Time)error
}
