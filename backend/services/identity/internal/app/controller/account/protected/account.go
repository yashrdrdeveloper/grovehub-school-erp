package protected
import "net/http"
// @Routes /accounts
type AccountController interface {
 // @POST /
 Create(http.ResponseWriter,*http.Request)
 // @POST /search
 Search(http.ResponseWriter,*http.Request)
 // @GET /{accountId}
 Get(http.ResponseWriter,*http.Request)
 // @PATCH /{accountId}
 Update(http.ResponseWriter,*http.Request)
 // @PATCH /{accountId}/suspend
 Suspend(http.ResponseWriter,*http.Request)
 // @PATCH /{accountId}/activate
 Activate(http.ResponseWriter,*http.Request)
 // @PATCH /{accountId}/deactivate
 Deactivate(http.ResponseWriter,*http.Request)
 // @POST /{accountId}/resend-invitation
 Resend(http.ResponseWriter,*http.Request)
}
