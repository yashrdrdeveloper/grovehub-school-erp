package controller

import "net/http"

type AuthController interface {
	Register(writer http.ResponseWriter, request *http.Request)
}
