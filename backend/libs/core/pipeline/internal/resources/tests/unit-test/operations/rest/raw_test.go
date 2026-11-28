package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"grovehub/libs/core/pipeline"
)

func TestRawWriter(t *testing.T) {
	writer := pipeline.RawWriter[string]()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, status := range []int{0, http.StatusOK, http.StatusCreated, http.StatusNoContent} {
		w := httptest.NewRecorder()
		result := pipeline.Result[string]{StatusCode: status, Data: "value", Headers: http.Header{"X-Test": []string{"one", "two"}}}
		if err := writer.Success(w, request, result); err != nil {
			t.Fatal(err)
		}
		expected := status
		if expected == 0 {
			expected = http.StatusOK
		}
		if w.Code != expected || len(w.Header().Values("X-Test")) != 2 {
			t.Fatalf("status=%d headers=%v", w.Code, w.Header())
		}
		if expected == http.StatusNoContent {
			if w.Body.Len() != 0 {
				t.Fatal("unexpected body")
			}
		} else if w.Body.String() != "\"value\"\n" || w.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("invalid raw JSON response: %q", w.Body.String())
		}
	}
	w := httptest.NewRecorder()
	if err := writer.Error(w, request, pipeline.Failure{StatusCode: http.StatusForbidden}); err != nil || w.Code != http.StatusForbidden || w.Body.Len() != 0 {
		t.Fatalf("invalid error response: %d %q %v", w.Code, w.Body.String(), err)
	}
	if pipeline.NoContent().StatusCode != http.StatusNoContent {
		t.Fatal("invalid empty response")
	}
}
