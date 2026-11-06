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
 Watch(http.ResponseWriter,*http.Request)
 GmailPush(http.ResponseWriter,*http.Request)
 MailProfile(http.ResponseWriter,*http.Request)
 MailMessages(http.ResponseWriter,*http.Request)
 MailMessage(http.ResponseWriter,*http.Request)
 MailHistory(http.ResponseWriter,*http.Request)
 MailSend(http.ResponseWriter,*http.Request)
 ExecuteAction(http.ResponseWriter,*http.Request)
}
