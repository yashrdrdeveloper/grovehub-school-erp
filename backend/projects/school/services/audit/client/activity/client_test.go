package activity

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/client/authentication"
)

func testClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	origin, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	httpClient := server.Client()
	httpClient.Transport = auth.ProofTransport{Base: httpClient.Transport}
	return New(origin, httpClient)
}

func TestSearch(t *testing.T) {
	const response = `{"success":true,"data":{"items":[],"pageInfo":{"rows":0,"hasPrevious":false,"hasNext":false}}}`
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/audit/restricted/search" {
			t.Errorf("unexpected operation: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-GH-Identity-Proof") != "test-proof" {
			t.Error("missing route-selected proof")
		}
		if r.Header.Get("Cookie") != "" {
			t.Error("browser cookies must not be forwarded")
		}
		var input struct {
			Scope string `json:"scope"`
			PageInfo struct { RowSize int `json:"rowSize"` } `json:"pageInfo"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Scope != "ME" || input.PageInfo.RowSize != 25 {
			t.Errorf("invalid generated request: %+v, %v", input, err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, response)
	})
	result, err := client.Search(context.Background(), json.RawMessage(`{"scope":"ME","pageInfo":{"rowSize":25}}`), "test-proof")
	if err != nil || result.Status != 200 || string(result.Body) != response {
		t.Fatalf("unexpected result: %+v, %v", result, err)
	}
}

func TestSearchRejectsInvalidInputBeforeNetwork(t *testing.T) {
	var calls atomic.Int32
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	})
	for _, input := range []string{"", "null", "[]", `{"unknown":true}`, "{} {}"} {
		_, err := client.Search(context.Background(), json.RawMessage(input), "test-proof")
		if !errors.Is(err, ErrInvalidRequest) {
			t.Errorf("expected invalid request for %q, got %v", input, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatal("invalid inputs reached Platform")
	}
}

func TestSearchPreservesForbiddenResponse(t *testing.T) {
	const response = `{"success":false,"errorCode":"FORBIDDEN","errorMessage":"denied","httpErrorCode":403,"errorMetaData":{}}`
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, response)
	})
	result, err := client.Search(context.Background(), json.RawMessage(`{"scope":"ME"}`), "test-proof")
	if err != nil || result.Status != http.StatusForbidden || string(result.Body) != response {
		t.Fatalf("unexpected result: %+v, %v", result, err)
	}
}

func TestSearchRejectsOversizedResponse(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "\"" + strings.Repeat("x", maxResponseBytes) + "\"")
	})
	if _, err := client.Search(context.Background(), json.RawMessage(`{"scope":"ME"}`), "test-proof"); err == nil {
		t.Fatal("oversized response accepted")
	}
}
