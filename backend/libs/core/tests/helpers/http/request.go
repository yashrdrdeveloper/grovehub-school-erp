// Package httphelpers provides in-memory HTTP request fixtures.
// It intentionally uses no external listener or network connection.
package httphelpers

import (
    "bytes"
    "net/http"
    "net/http/httptest"
)

func Request(method,url string,body []byte) (*http.Request,*httptest.ResponseRecorder) {
    request:=httptest.NewRequest(method,url,bytes.NewReader(body))
    if body!=nil {request.Header.Set("Content-Type","application/json")}
    return request,httptest.NewRecorder()
}
