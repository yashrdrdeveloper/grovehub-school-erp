package controller

import "net/http"

type ConnectionController interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	Disconnect(http.ResponseWriter, *http.Request)
	BeginOAuth(http.ResponseWriter, *http.Request)
	FinishOAuth(http.ResponseWriter, *http.Request)
	RunTest(http.ResponseWriter, *http.Request)
	ListTests(http.ResponseWriter, *http.Request)
	Watch(http.ResponseWriter, *http.Request)
	GmailPush(http.ResponseWriter, *http.Request)
	ExecuteAction(http.ResponseWriter, *http.Request)
}
