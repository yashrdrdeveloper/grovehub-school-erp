package public

import "net/http"

// @Routes /public/auth
type AuthController interface {
    // @POST /login
    Login(writer http.ResponseWriter, request *http.Request)
}
