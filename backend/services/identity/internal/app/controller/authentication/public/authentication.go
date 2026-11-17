package public
import "net/http"
// @Routes /auth
type AuthenticationController interface{
 // @POST /login
 // @RouteLimit 12
 Login(http.ResponseWriter,*http.Request)
 // @POST /refresh
 // @RouteLimit 30
 Refresh(http.ResponseWriter,*http.Request)
 // @POST /forgot-password
 // @RouteLimit 5
 ForgotPassword(http.ResponseWriter,*http.Request)
 // @POST /recovery/exchange
 // @RouteLimit 8
 RecoveryExchange(http.ResponseWriter,*http.Request)
 // @POST /reset-password
 ResetPassword(http.ResponseWriter,*http.Request)
 // @POST /invitations/exchange
 // @RouteLimit 8
 InvitationExchange(http.ResponseWriter,*http.Request)
}
