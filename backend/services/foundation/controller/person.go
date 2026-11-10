package controller

import "net/http"

// @Routes /foundation/persons
type PersonController interface {
    // @POST /
    Create(writer http.ResponseWriter, request *http.Request)
    // @GET /{gid}
    GetByGID(writer http.ResponseWriter, request *http.Request)
    // @PATCH /{gid}
    Update(writer http.ResponseWriter, request *http.Request)
    // @POST /search
    Search(writer http.ResponseWriter, request *http.Request)
    // @PATCH /{gid}/deactivate
    Deactivate(writer http.ResponseWriter, request *http.Request)
}
