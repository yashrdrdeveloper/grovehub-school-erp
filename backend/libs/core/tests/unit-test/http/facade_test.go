package corehttp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
)

func TestFacade(t *testing.T) {
	if corehttp.RequestID != "X-Request-ID" || corehttp.IdempotencyKey != "Idempotency-Key" { t.Fatal("header names changed") }
	if corehttp.ID(corehttp.Ensure(context.Background())) == "" { t.Fatal("missing correlation") }
	var drain atomic.Bool
	handler := corehttp.ReadinessHandler(&drain, func(context.Context) error { return nil })
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/ready", nil))
	if w.Code != 204 { t.Fatalf("readiness status %d", w.Code) }
	var body map[string]string
	if err := corehttp.DecodeJSON(httptest.NewRequest("POST","/",strings.NewReader(`{"name":"a"}`)), &body); err != nil { t.Fatal(err) }
	value, status, err := corehttp.Call[string](context.Background(), http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if err := corehttp.WriteSuccess(w,200,body["name"]); err != nil { t.Error(err) }
	}), "GET","/",nil,nil)
	if err != nil || status != 200 || value != "a" { t.Fatalf("unexpected response %v %d %v",value,status,err) }
}
