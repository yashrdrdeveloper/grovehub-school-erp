package controller

import "net/http"

type ActionController interface{ExecuteAction(http.ResponseWriter,*http.Request)}
