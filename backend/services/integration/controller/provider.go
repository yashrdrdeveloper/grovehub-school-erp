package controller

import "net/http"

type ProviderController interface{
 Providers(http.ResponseWriter,*http.Request)
 Provider(http.ResponseWriter,*http.Request)
}
