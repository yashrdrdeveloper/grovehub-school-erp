package public
import "net/http"
// @Routes /auth
type AuthenticationController interface{
 // @POST /login
 Login(http.ResponseWriter,*http.Request)
 // @POST /refresh
 Refresh(http.ResponseWriter,*http.Request)
 // @POST /forgot-password
 ForgotPassword(http.ResponseWriter,*http.Request)
 // @POST /recovery/exchange
 RecoveryExchange(http.ResponseWriter,*http.Request)
 // @POST /reset-password
 ResetPassword(http.ResponseWriter,*http.Request)
 // @POST /invitations/exchange
 InvitationExchange(http.ResponseWriter,*http.Request)
}
