package controller

import "net/http"

type SenderController interface {
	Register(writer http.ResponseWriter, request *http.Request)
	ListAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Deactivate(writer http.ResponseWriter, request *http.Request)
}
