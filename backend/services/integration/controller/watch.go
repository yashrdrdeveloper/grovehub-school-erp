package controller

import "net/http"

type WatchController interface{Watch(http.ResponseWriter,*http.Request)}
