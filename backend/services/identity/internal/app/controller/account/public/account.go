package public
import "net/http"
// @Routes /identity/public/accounts
type Controller interface{Setup(http.ResponseWriter,*http.Request)}
