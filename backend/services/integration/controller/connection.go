package controller

import "net/http"

type ConnectionController interface{
 Create(http.ResponseWriter,*http.Request)
 Get(http.ResponseWriter,*http.Request)
 List(http.ResponseWriter,*http.Request)
 Disconnect(http.ResponseWriter,*http.Request)
}
