package controller

import "net/http"

type DomainController interface {
	Register(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	GetByName(writer http.ResponseWriter, request *http.Request)
	Verify(writer http.ResponseWriter, request *http.Request)
	Deactivate(writer http.ResponseWriter, request *http.Request)
}
