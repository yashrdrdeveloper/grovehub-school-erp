package protected
import "net/http"
// @Routes /identity/protected/auth
type Controller interface{
 Logout(http.ResponseWriter,*http.Request)
 LogoutAll(http.ResponseWriter,*http.Request)
 ChangePassword(http.ResponseWriter,*http.Request)
 Reauthenticate(http.ResponseWriter,*http.Request)
 Sessions(http.ResponseWriter,*http.Request)
 RevokeSession(http.ResponseWriter,*http.Request)
 Me(http.ResponseWriter,*http.Request)
}
