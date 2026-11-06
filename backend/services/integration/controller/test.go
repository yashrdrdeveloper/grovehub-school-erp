package controller

import "net/http"

type TestController interface{
 RunTest(http.ResponseWriter,*http.Request)
 ListTests(http.ResponseWriter,*http.Request)
}
