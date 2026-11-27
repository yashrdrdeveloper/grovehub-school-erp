package authentication

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGeneratedLogin(t *testing.T) {
	const response = `{"success":true,"data":{"expiresIn":900,"requiresSetup":false}}`
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/identity/restricted/auth/login" {
			t.Errorf("unexpected operation: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-GH-Identity-Proof") != "" || r.Header.Get("Cookie") != "" {
			t.Error("public login forwarded an identity credential")
		}
		var input Login
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Identifier != "test@example.com" || input.Password != "test-password" {
			t.Error("invalid generated login request")
		}
		w.Header().Set("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{Name: "__Host-gh_access", Value: "opaque-test", Path: "/", Secure: true, HttpOnly: true})
		_, _ = io.WriteString(w, response)
	})
	result, err := client.Login(WithIdentityProof(context.Background(), "must-not-forward"), Login{Identifier:"test@example.com", Password:"test-password"})
	if err != nil || result.Status != http.StatusOK || string(result.Body) != response || len(result.Cookies) != 1 || result.Cookies[0].Name != "__Host-gh_access" || !result.Cookies[0].HttpOnly {
		t.Fatalf("login response or cookie not preserved: %v", err)
	}
}

func TestGeneratedForgotPassword(t *testing.T) {
	const response = `{"success":true,"data":{"accepted":true,"delivery":"EMAIL"}}`
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/identity/restricted/auth/forgot-password" {
			t.Errorf("unexpected operation: %s %s", r.Method, r.URL.Path)
		}
		var input ForgotPassword
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Email != "test@example.com" {
			t.Error("invalid generated recovery request")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, response)
	})
	result, err := client.ForgotPassword(context.Background(), ForgotPassword{Email:"test@example.com"})
	if err != nil || result.Status != http.StatusAccepted || string(result.Body) != response {
		t.Fatalf("recovery response not preserved: %v", err)
	}
}

func TestGeneratedLoginPreservesRejection(t *testing.T) {
	const response = `{"success":false,"errorCode":"UNAUTHORIZED","errorMessage":"denied","httpErrorCode":401,"errorMetaData":{}}`
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, response)
	})
	result, err := client.Login(context.Background(), Login{Identifier:"test@example.com", Password:"test-password"})
	if err != nil || result.Status != http.StatusUnauthorized || string(result.Body) != response {
		t.Fatalf("login rejection not preserved: %v", err)
	}
}

func TestGeneratedLoginRejectsInvalidResponses(t *testing.T) {
	for _, body := range []string{"not-json", "\"" + strings.Repeat("x", 1<<20) + "\""} {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, body)
		})
		if _, err := client.Login(context.Background(), Login{Identifier:"test@example.com", Password:"test-password"}); err == nil {
			t.Fatal("invalid platform response accepted")
		}
	}
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	origin, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	client := server.Client()
	client.Transport = ProofTransport{Base:client.Transport}
	return New(origin, client)
}
