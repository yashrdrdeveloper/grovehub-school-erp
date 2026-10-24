package controller

import "net/http"

type PersonController interface {
	Create(writer http.ResponseWriter, request *http.Request)
	GetByGID(writer http.ResponseWriter, request *http.Request)
	Update(writer http.ResponseWriter, request *http.Request)
	Search(writer http.ResponseWriter, request *http.Request)
	Deactivate(writer http.ResponseWriter, request *http.Request)
}
