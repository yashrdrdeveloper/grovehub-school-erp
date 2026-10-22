package routing_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	responsedto "grovehub/libs/core/dto/response"
	coreerrors "grovehub/libs/core/errors"
	"grovehub/launcher/routing"
)

func TestFallbackReturnsStandardNotFoundResponse(t *testing.T) {
	mux := http.NewServeMux()
	routing.RegisterFallback(mux)

	request := httptest.NewRequest(http.MethodGet, "/abc", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	var response responsedto.Error[struct{}]
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Success {
		t.Fatal("expected success to be false")
	}
	if response.ErrorCode != coreerrors.Code("PLATFORM_ROUTE_NOT_FOUND") {
		t.Fatalf("unexpected error code: %s", response.ErrorCode)
	}
	if response.ErrorMessage != "API route not found" {
		t.Fatalf("unexpected error message: %s", response.ErrorMessage)
	}
	if response.HTTPErrorCode != http.StatusNotFound {
		t.Fatalf("expected http error code %d, got %d", http.StatusNotFound, response.HTTPErrorCode)
	}
}

func TestFallbackDoesNotOverrideRegisteredRoute(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /identity/health", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})
	routing.RegisterFallback(mux)

	request := httptest.NewRequest(http.MethodGet, "/identity/health", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
}
