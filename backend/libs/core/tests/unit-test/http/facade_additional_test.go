package corehttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	corehttp "grovehub/libs/core/http"
)

func TestFacadeEnvelopesAndStrictDecoding(t *testing.T) {
	success := corehttp.NewSuccess("data")
	if !success.Success || success.Data != "data" { t.Fatal(success) }
	failure := corehttp.NewError("INVALID", "invalid request", http.StatusBadRequest, "field")
	if failure.Success || failure.ErrorCode != "INVALID" || failure.HTTPErrorCode != 400 || failure.ErrorMetaData != "field" { t.Fatal(failure) }
	var body struct { Name string `json:"name"` }
	if err := corehttp.DecodeStrictJSON(httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Ada"}`)), &body); err != nil || body.Name != "Ada" { t.Fatalf("body=%+v err=%v",body,err) }
	if err := corehttp.DecodeStrictJSON(httptest.NewRequest("POST", "/", strings.NewReader(`{"unknown":1}`)), &body); err == nil { t.Fatal("unknown field accepted") }
}

func TestFacadeHTTPClients(t *testing.T) {
	for _, client := range []*http.Client{corehttp.NewClient(time.Second), corehttp.NewServiceClient("identity",time.Second)} {
		if client.Timeout != time.Second || client.Transport == nil { t.Fatal("client configuration lost") }
	}
	if corehttp.NewServiceTransport("identity",http.DefaultTransport) == nil { t.Fatal("missing service transport") }
}
