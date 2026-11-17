package protected
import "net/http"
// @Routes /auth
type AuthenticationController interface{
 // @POST /logout
 Logout(http.ResponseWriter,*http.Request)
 // @POST /logout-all
 LogoutAll(http.ResponseWriter,*http.Request)
 // @POST /change-password
 ChangePassword(http.ResponseWriter,*http.Request)
 // @POST /reauthenticate
 Reauthenticate(http.ResponseWriter,*http.Request)
 // @GET /sessions
 Sessions(http.ResponseWriter,*http.Request)
 // @DELETE /sessions/{sessionId}
 RevokeSession(http.ResponseWriter,*http.Request)
 // @GET /me
 Me(http.ResponseWriter,*http.Request)
}
