package controller

import "net/http"

type ConnectionController interface {
 Providers(http.ResponseWriter,*http.Request)
 Create(http.ResponseWriter,*http.Request)
 Get(http.ResponseWriter,*http.Request)
 List(http.ResponseWriter,*http.Request)
 Disconnect(http.ResponseWriter,*http.Request)
 BeginOAuth(http.ResponseWriter,*http.Request)
 FinishOAuth(http.ResponseWriter,*http.Request)
 RunTest(http.ResponseWriter,*http.Request)
 ListTests(http.ResponseWriter,*http.Request)
 Sync(http.ResponseWriter,*http.Request)
}
