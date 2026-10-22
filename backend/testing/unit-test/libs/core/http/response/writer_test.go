package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	responsedto "grovehub/libs/core/dto/response"
	coreerrors "grovehub/libs/core/errors"
	httpresponse "grovehub/libs/core/http/response"
)

func TestWriteErrorUsesDefinitionAndHTTPStatus(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	definition := coreerrors.NewDefinition(
		"IDENTITY_SERVICE_DOWN",
		"Identity service is unavailable",
	)

	if err := httpresponse.WriteError(recorder, http.StatusServiceUnavailable, definition); err != nil {
		t.Fatalf("write error response: %v", err)
	}

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, recorder.Code)
	}

	var payload responsedto.Error[map[string]any]
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error response: %v", err)
	}

	if payload.Success {
		t.Fatal("expected success to be false")
	}
	if payload.ErrorCode != definition.Code() {
		t.Fatalf("expected error code %q, got %q", definition.Code(), payload.ErrorCode)
	}
	if payload.ErrorMessage != definition.Message() {
		t.Fatalf("expected error message %q, got %q", definition.Message(), payload.ErrorMessage)
	}
	if payload.HTTPErrorCode != http.StatusServiceUnavailable {
		t.Fatalf("expected httpErrorCode %d, got %d", http.StatusServiceUnavailable, payload.HTTPErrorCode)
	}
	if len(payload.ErrorMetaData) != 0 {
		t.Fatalf("expected empty metadata, got %#v", payload.ErrorMetaData)
	}
}

func TestWriteErrorWithMetadata(t *testing.T) {
	t.Parallel()

	type metadata struct {
		Field string `json:"field"`
	}

	recorder := httptest.NewRecorder()
	definition := coreerrors.NewDefinition("IDENTITY_INVALID_REQUEST", "Identity request is invalid")

	if err := httpresponse.WriteErrorWithMetadata(
		recorder,
		http.StatusBadRequest,
		definition,
		metadata{Field: "email"},
	); err != nil {
		t.Fatalf("write error response with metadata: %v", err)
	}

	var payload responsedto.Error[metadata]
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error response: %v", err)
	}

	if payload.ErrorMetaData.Field != "email" {
		t.Fatalf("expected metadata field %q, got %q", "email", payload.ErrorMetaData.Field)
	}
}
