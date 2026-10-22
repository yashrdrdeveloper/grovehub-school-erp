package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	responsedto "grovehub/libs/core/dto/response"
	"grovehub/services/identity/controller"
	identityresponse "grovehub/services/identity/dto/v1/response"
)

func TestHealthControllerReturnsStandardSuccessEnvelope(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/identity/health", nil)

	controller.NewHealthController().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var payload responsedto.Success[identityresponse.HealthData]
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode health response: %v", err)
	}

	if !payload.Success {
		t.Fatal("expected success to be true")
	}
	if payload.Data.Service != "identity" {
		t.Fatalf("expected service %q, got %q", "identity", payload.Data.Service)
	}
	if payload.Data.Status != "ok" {
		t.Fatalf("expected status %q, got %q", "ok", payload.Data.Status)
	}
}
