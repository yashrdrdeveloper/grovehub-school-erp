package gateway_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	ginkgo "github.com/onsi/ginkgo/v2"

	external "grovehub/launcher/internal/src/services/gateway/operations/external/impl"
	schoolgateway "grovehub/libs/platform/session"
)

type accessVerifier struct{}

func (accessVerifier) VerifyAccess(_ context.Context, secret string) (schoolgateway.Principal, error) {
	if secret != "test-secret" {
		return schoolgateway.Principal{}, context.Canceled
	}
	return schoolgateway.Principal{
		AccountID: "account",
		PersonGID: "person",
		TenantID:  "tenant",
		SessionID: "session",
	}, nil
}

func verifyProtectedWildcardRoutingPreservesPathValues() {
	t := ginkgo.GinkgoT()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /foundation/protected/persons/{gid}", func(w http.ResponseWriter, r *http.Request) {
		principal, ok := schoolgateway.PrincipalFrom(r.Context())
		if !ok || principal.AccountID != "account" || r.PathValue("gid") != "p-123" {
			t.Error("verified principal or wildcard binding was lost")
		}

		w.WriteHeader(http.StatusOK)
	})

	gateway := external.New(mux, accessVerifier{}, nil, "https://app.example.test")
	request := httptest.NewRequest(http.MethodGet, "/foundation/protected/persons/p-123", nil)
	request.AddCookie(&http.Cookie{Name: schoolgateway.AccessCookie, Value: "test-secret"})
	recorder := httptest.NewRecorder()
	gateway.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

func verifyUnknownAndRestrictedRoutesFailClosed() {
	t := ginkgo.GinkgoT()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /unclassified", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /integrations/restricted/webhooks/gmail", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	gateway := external.New(mux, accessVerifier{}, nil, "https://app.example.test")
	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/unclassified"},
		{http.MethodPost, "/integrations/restricted/webhooks/gmail"},
	} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(tc.method, tc.path, nil)
		gateway.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("%s: expected forbidden, got %d", tc.path, recorder.Code)
		}
	}
}

func verifyPublicUnsafeRequestsRequireTrustedOrigin() {
	t := ginkgo.GinkgoT()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /identity/public/auth/login", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	gateway := external.New(mux, accessVerifier{}, nil, "https://app.example.test")
	request := httptest.NewRequest(http.MethodPost, "/identity/public/auth/login", nil)
	request.Header.Set("Origin", "https://attacker.example")
	recorder := httptest.NewRecorder()
	gateway.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("cross-origin login allowed: %d", recorder.Code)
	}
}

func verifyRateLimitUsesProvidedExactMethodAndRoute() {
	t := ginkgo.GinkgoT()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /identity/public/auth/login", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	limits := map[string]int{"POST /identity/public/auth/login": 2}
	gateway := external.New(mux, accessVerifier{}, nil, "https://app.example.test", limits)
	for attempt := 0; attempt < 3; attempt++ {
		request := httptest.NewRequest(http.MethodPost, "/identity/public/auth/login", nil)
		request.Header.Set("Origin", "https://app.example.test")
		recorder := httptest.NewRecorder()
		gateway.ServeHTTP(recorder, request)
		want := http.StatusOK
		if attempt == 2 {
			want = http.StatusTooManyRequests
		}

		if recorder.Code != want {
			t.Fatalf("attempt %d: expected %d, got %d", attempt+1, want, recorder.Code)
		}
	}
}

func verifyLocalViteLoopbackOriginsShareSamePort() {
	t := ginkgo.GinkgoT()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /identity/public/auth/invitations/exchange", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	gateway := external.New(mux, accessVerifier{}, nil, "http://127.0.0.1:4173")
	cases := []struct {
		origin string
		status int
	}{
		{"http://localhost:4173", http.StatusOK},
		{"http://127.0.0.1:4173", http.StatusOK},
		{"http://localhost:4174", http.StatusForbidden},
		{"https://localhost:4173", http.StatusForbidden},
		{"http://evil.example:4173", http.StatusForbidden},
		{"", http.StatusForbidden},
	}
	for _, tc := range cases {
		request := httptest.NewRequest(http.MethodPost, "/identity/public/auth/invitations/exchange", nil)
		if tc.origin != "" {
			request.Header.Set("Origin", tc.origin)
		}

		recorder := httptest.NewRecorder()
		gateway.ServeHTTP(recorder, request)
		if recorder.Code != tc.status {
			t.Errorf("origin %q: expected %d, got %d", tc.origin, tc.status, recorder.Code)
		}
	}
}

type routeAccess struct {
	accessVerifier
	allowed bool
	pattern string
}

func (access *routeAccess) AuthorizeProtectedRoute(ctx context.Context, pattern string) error {
	access.pattern = pattern
	if _, ok := schoolgateway.PrincipalFrom(ctx); !ok || !access.allowed {
		return context.Canceled
	}

	return nil
}

func verifyIdentityProtectedRoutePolicy() {
	t := ginkgo.GinkgoT()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /identity/protected/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	access := &routeAccess{}
	for _, tc := range []struct {
		verifier schoolgateway.AccessVerifier
		want     int
	}{
		{accessVerifier{}, http.StatusForbidden},
		{access, http.StatusForbidden},
		{&routeAccess{allowed: true}, http.StatusOK},
	} {
		request := httptest.NewRequest(http.MethodGet, "/identity/protected/accounts/target", nil)
		request.AddCookie(&http.Cookie{Name: schoolgateway.AccessCookie, Value: "test-secret"})
		recorder := httptest.NewRecorder()
		external.New(mux, tc.verifier, nil, "").ServeHTTP(recorder, request)
		if recorder.Code != tc.want {
			t.Fatalf("expected %d, got %d", tc.want, recorder.Code)
		}
	}

	if access.pattern != "GET /identity/protected/accounts/{accountId}" {
		t.Fatalf("expected registered pattern, got %q", access.pattern)
	}
}

// @Testing suite=LauncherExternalGateway
func ExternalGatewaySpecs() {
	ginkgo.Describe("ExternalGateway", func() {
		ginkgo.It("Protected Wildcard Routing Preserves Path Values", verifyProtectedWildcardRoutingPreservesPathValues)
		ginkgo.It("Unknown And Restricted Routes Fail Closed", verifyUnknownAndRestrictedRoutesFailClosed)
		ginkgo.It("Public Unsafe Requests Require Trusted Origin", verifyPublicUnsafeRequestsRequireTrustedOrigin)
		ginkgo.It("Rate Limit Uses Provided Exact Method And Route", verifyRateLimitUsesProvidedExactMethodAndRoute)
		ginkgo.It("Local Vite Loopback Origins Share Same Port", verifyLocalViteLoopbackOriginsShareSamePort)
		ginkgo.It("Identity Protected Route Policy", verifyIdentityProtectedRoutePolicy)
	})
}
