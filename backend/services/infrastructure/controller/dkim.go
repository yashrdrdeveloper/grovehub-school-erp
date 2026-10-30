package controller

import "net/http"

type DKIMController interface {
	Ensure(writer http.ResponseWriter, request *http.Request)
}
