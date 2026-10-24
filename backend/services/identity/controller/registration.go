package controller

import "net/http"

type RegistrationController interface {
	Register(writer http.ResponseWriter, request *http.Request)
}
