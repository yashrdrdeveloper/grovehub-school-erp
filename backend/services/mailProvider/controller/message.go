package controller

import "net/http"

type MessageController interface {
	Send(writer http.ResponseWriter, request *http.Request)
	Search(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
}
