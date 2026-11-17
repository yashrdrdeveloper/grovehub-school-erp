package controller

import "net/http"

// Legacy controller contract; the routed contract lives in controller/protected.
type AuthController interface {
    // @POST /register
    Register(writer http.ResponseWriter, request *http.Request)
}
