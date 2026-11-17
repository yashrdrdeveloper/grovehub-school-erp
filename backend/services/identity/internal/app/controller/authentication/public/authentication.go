package public
import "net/http"
// @Routes /identity/public/auth
type Controller interface{
 Login(http.ResponseWriter,*http.Request)
 Refresh(http.ResponseWriter,*http.Request)
 ForgotPassword(http.ResponseWriter,*http.Request)
 RecoveryExchange(http.ResponseWriter,*http.Request)
 ResetPassword(http.ResponseWriter,*http.Request)
 InvitationExchange(http.ResponseWriter,*http.Request)
}
