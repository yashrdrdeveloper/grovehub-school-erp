package protected
import "net/http"
// @Routes /identity/protected/accounts
type Controller interface {
 Create(http.ResponseWriter,*http.Request)
 Search(http.ResponseWriter,*http.Request)
 Get(http.ResponseWriter,*http.Request)
 Update(http.ResponseWriter,*http.Request)
 Suspend(http.ResponseWriter,*http.Request)
 Activate(http.ResponseWriter,*http.Request)
 Deactivate(http.ResponseWriter,*http.Request)
 Resend(http.ResponseWriter,*http.Request)
}
