package controller

import "net/http"

// @Routes /protected/auth
type AuthController interface {
    // @POST /register
    Register(writer http.ResponseWriter, request *http.Request)
}
