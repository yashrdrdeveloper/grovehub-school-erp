package public
import "net/http"
// @Routes /accounts
type AccountController interface{// @POST /setup
Setup(http.ResponseWriter,*http.Request)}
