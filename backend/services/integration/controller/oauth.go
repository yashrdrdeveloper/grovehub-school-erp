package controller

import "net/http"

type OAuthController interface{
 BeginOAuth(http.ResponseWriter,*http.Request)
 FinishOAuth(http.ResponseWriter,*http.Request)
}
