package workload_test

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	fixture "grovehub/launcher/internal/resources/tests/helpers/workload"
	mockexternal "grovehub/launcher/internal/resources/tests/mocks/gateway/external"
	mock "grovehub/launcher/internal/resources/tests/mocks/gateway/workload"
	"grovehub/launcher/internal/src/services/gateway/operations/workload"
	"grovehub/launcher/internal/src/services/gateway/policies"
	core "grovehub/libs/core/gateway"
	"grovehub/libs/platform/session"
	mail "grovehub/services/integration/client/mail"
	integration "grovehub/services/integration/client/principal"
)

var delegated *mock.MockDelegatedIdentity

// @BeforeEach
func SetupDelegatedIdentity() {
	delegated = mock.NewMockDelegatedIdentity(gomock.NewController(GinkgoT()))
}

func workloadJSON(config workload.WorkloadConfig) string {
	raw, err := json.Marshal(config)
	Expect(err).NotTo(HaveOccurred())
	return string(raw)
}

func workloadRoutes(patterns ...string) map[string]core.RoutePolicy {
	routes := make(map[string]core.RoutePolicy, len(patterns))
	for _, pattern := range patterns {
		policy, ok := policies.Find(pattern)
		Expect(ok).To(BeTrue(), pattern)
		routes[pattern] = policy
	}

	return routes
}

func workloadScopes(routes map[string]core.RoutePolicy) []string {
	scopes := make([]string, 0, len(routes))
	for _, policy := range routes {
		scopes = append(scopes, policy.Scope)
	}

	return scopes
}

func signedCall(verifier core.WorkloadVerifier, pattern, path string, key core.WorkloadKeyConfig, private ed25519.PrivateKey, claims core.WorkloadClaims, headers map[string][]string) (context.Context, error) {
	request, err := fixture.Request(pattern, path, key, private, claims)
	Expect(err).NotTo(HaveOccurred())
	for header, values := range headers {
		for _, value := range values {
			request.Header.Add(header, value)
		}
	}

	return verifier.VerifyWorkload(context.Background(), request)
}

// @Testing suite=LauncherWorkloads
func GenericWorkloadSpecs() {
	Describe("NewPlatformWorkloads", func() {
		It("rejects malformed, trailing and incomplete configuration and unsafe route grants", func() {
			for _, raw := range []string{"{", "{\"unexpected\":true}", "{} {}", "{}"} {
				_, err := workload.NewPlatformWorkloads(raw, nil, nil)
				Expect(err).To(HaveOccurred())
			}

			key, _, err := fixture.Key("school", "school", "school", []string{"identity.auth.me"})
			Expect(err).NotTo(HaveOccurred())
			for _, cfg := range []workload.WorkloadConfig{
				{Keys: []core.WorkloadKeyConfig{key}},
				{Keys: []core.WorkloadKeyConfig{key}, Routes: map[string]core.RoutePolicy{"GET /unknown": {Scope: "identity.auth.me"}}},
				{Keys: []core.WorkloadKeyConfig{key}, Routes: workloadRoutes("GET /identity/restricted/auth/me")},
			} {
				_, err := workload.NewPlatformWorkloads(workloadJSON(cfg), nil, nil)
				Expect(err).To(HaveOccurred())
			}

			key.Scopes = []string{mail.ReconciliationScope}
			_, err = workload.NewPlatformWorkloads(workloadJSON(workload.WorkloadConfig{Keys: []core.WorkloadKeyConfig{key}, Routes: workloadRoutes("POST " + mail.ReconciliationPath)}), nil, nil)
			Expect(err).To(HaveOccurred())
		})
		It("delegates only the exact Gmail webhook route to its isolated verifier", func() {
			request := httptest.NewRequest(http.MethodPost, "/integrations/restricted/webhooks/gmail", nil)
			request.Pattern = "POST /integrations/restricted/webhooks/gmail"
			empty, err := workload.NewPlatformWorkloads("", nil, nil)
			Expect(err).NotTo(HaveOccurred())
			_, err = empty.VerifyWorkload(context.Background(), nil)
			Expect(err).To(MatchError(core.ErrWorkloadDenied))
			_, err = empty.VerifyWorkload(context.Background(), request)
			Expect(err).To(MatchError(core.ErrWorkloadDenied))
			webhook := mockexternal.NewMockWorkloadVerifier(gomock.NewController(GinkgoT()))
			ctx := request.Context()
			webhook.EXPECT().VerifyWorkload(ctx, request).Return(ctx, nil)
			v, err := workload.NewPlatformWorkloads("", nil, webhook)
			Expect(err).NotTo(HaveOccurred())
			got, err := v.VerifyWorkload(ctx, request)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(ctx))
		})
	})
	Describe("Generic route authorization", func() {
		It("requires separate Foundation, School and reconciliation workload identities", func() {
			foundationPattern := "POST /foundation/restricted/persons"
			loginPattern := "POST /identity/restricted/auth/login"
			mailPattern := "POST " + mail.ReconciliationPath
			routes := workloadRoutes(foundationPattern, loginPattern, mailPattern)
			foundation, foundationPrivate, err := fixture.Key("foundation", "school-foundation", "school-foundation", []string{routes[foundationPattern].Scope})
			Expect(err).NotTo(HaveOccurred())
			school, schoolPrivate, err := fixture.Key("school", "school", "school", []string{routes[foundationPattern].Scope, routes[loginPattern].Scope})
			Expect(err).NotTo(HaveOccurred())
			school.AllowDelegation = true
			operator, operatorPrivate, err := fixture.Key("mail", mail.ReconciliationSubject, mail.ReconciliationSubject, []string{mail.ReconciliationScope})
			Expect(err).NotTo(HaveOccurred())
			v, err := workload.NewPlatformWorkloads(workloadJSON(workload.WorkloadConfig{Keys: []core.WorkloadKeyConfig{foundation, school, operator}, Routes: routes}), nil, nil)
			Expect(err).NotTo(HaveOccurred())
			for _, tc := range []struct {
				pattern, path string
				key           core.WorkloadKeyConfig
				private       ed25519.PrivateKey
				scopes        []string
				sessionID     string
				headers       map[string][]string
				allowed       bool
			}{
				{foundationPattern, "/foundation/restricted/persons", foundation, foundationPrivate, foundation.Scopes, "", nil, true},
				{foundationPattern, "/foundation/restricted/persons", school, schoolPrivate, []string{routes[foundationPattern].Scope}, "", nil, false},
				{loginPattern, "/identity/restricted/auth/login", school, schoolPrivate, []string{routes[loginPattern].Scope}, "", nil, true},
				{loginPattern, "/identity/restricted/auth/login", school, schoolPrivate, []string{routes[loginPattern].Scope}, "asserted-session", nil, false},
				{loginPattern, "/identity/restricted/auth/login", school, schoolPrivate, []string{routes[loginPattern].Scope, routes[loginPattern].Scope}, "", nil, false},
				{mailPattern, mail.ReconciliationPath, operator, operatorPrivate, operator.Scopes, "", nil, true},
				{mailPattern, mail.ReconciliationPath + "?unexpected=1", operator, operatorPrivate, operator.Scopes, "", nil, false},
				{mailPattern, mail.ReconciliationPath, operator, operatorPrivate, operator.Scopes, "", map[string][]string{"Cookie": {"untrusted=1"}}, false},
				{mailPattern, mail.ReconciliationPath, operator, operatorPrivate, operator.Scopes, "", map[string][]string{"X-GH-Identity-Proof": {"untrusted"}}, false},
			} {
				_, err := signedCall(v, tc.pattern, tc.path, tc.key, tc.private, core.WorkloadClaims{Scopes: tc.scopes, SessionID: tc.sessionID}, tc.headers)
				if tc.allowed {
					Expect(err).NotTo(HaveOccurred(), tc.pattern)
				} else {
					Expect(err).To(MatchError(core.ErrWorkloadDenied), tc.pattern)
				}
			}
		})
		It("enforces actor proof and authorization on all Integration routes", func() {
			patterns := []string{
				"POST /integrations/restricted/connections",
				"GET /integrations/restricted/connections/{id}",
				"GET /integrations/restricted/connections",
				"DELETE /integrations/restricted/connections/{id}",
				"POST /integrations/restricted/connections/{id}/oauth/start",
				"POST /integrations/restricted/connections/{id}/test",
				"GET /integrations/restricted/connections/{id}/tests",
				"POST /integrations/restricted/connections/{id}/watch",
				"POST /integrations/restricted/connections/system-mail/sender",
				"POST /integrations/restricted/connections/{id}/actions/{action}",
				"GET /integrations/restricted/providers",
				"GET /integrations/restricted/providers/{provider}",
			}
			routes := workloadRoutes(patterns...)
			key, private, err := fixture.Key("school", "school", "school", workloadScopes(routes))
			Expect(err).NotTo(HaveOccurred())
			actor := session.Principal{AccountID: "account", SessionID: "session", TenantID: "tenant", PersonGID: "person", SessionGeneration: 1}
			delegated.EXPECT().VerifyAccess(gomock.Any(), "proof").Return(actor, nil).Times(len(patterns))
			v, err := workload.NewPlatformWorkloads(workloadJSON(workload.WorkloadConfig{Keys: []core.WorkloadKeyConfig{key}, Routes: routes}), delegated, nil)
			Expect(err).NotTo(HaveOccurred())
			for _, pattern := range patterns {
				_, path, _ := strings.Cut(pattern, " ")
				path = strings.NewReplacer("{id}", "connection", "{action}", "send", "{provider}", "gmail").Replace(path)
				ctx, err := signedCall(v, pattern, path, key, private, core.WorkloadClaims{Scopes: []string{routes[pattern].Scope}}, map[string][]string{"X-GH-Identity-Proof": {"proof"}})
				Expect(err).NotTo(HaveOccurred(), pattern)
				got, ok := integration.From(ctx)
				Expect(ok).To(BeTrue())
				Expect(got).To(Equal(integration.VerifiedPrincipal{TenantID: "tenant", ActorGID: "person"}))
			}

			delegated.EXPECT().VerifyAccess(gomock.Any(), "incomplete").Return(session.Principal{AccountID: "account", SessionID: "session"}, nil)
			_, err = signedCall(v, patterns[0], "/integrations/restricted/connections", key, private, core.WorkloadClaims{Scopes: []string{routes[patterns[0]].Scope}}, map[string][]string{"X-GH-Identity-Proof": {"incomplete"}})
			Expect(err).To(MatchError(core.ErrWorkloadDenied))
		})
		It("rejects invalid delegated proof, untrusted actor assertions, incomplete principals and denied protected routes", func() {
			pattern := "GET /identity/restricted/auth/me"
			path := "/identity/restricted/auth/me"
			routes := workloadRoutes(pattern)
			scope := routes[pattern].Scope
			key, private, err := fixture.Key("school", "school", "school", []string{scope})
			Expect(err).NotTo(HaveOccurred())
			key.AllowDelegation = true
			foreign, foreignPrivate, err := fixture.Key("foreign", "foreign", "foreign", []string{scope})
			Expect(err).NotTo(HaveOccurred())
			v, err := workload.NewPlatformWorkloads(workloadJSON(workload.WorkloadConfig{Keys: []core.WorkloadKeyConfig{key, foreign}, Routes: routes}), delegated, nil)
			Expect(err).NotTo(HaveOccurred())
			for _, proof := range [][]string{nil, {""}, {"a", "b"}, {strings.Repeat("a", 2049)}, {"invalid proof"}} {
				_, err := signedCall(v, pattern, path, key, private, core.WorkloadClaims{Scopes: []string{scope}}, map[string][]string{"X-GH-Identity-Proof": proof})
				Expect(err).To(MatchError(core.ErrWorkloadDenied))
			}

			_, err = signedCall(v, pattern, path, foreign, foreignPrivate, core.WorkloadClaims{Scopes: []string{scope}}, map[string][]string{"X-GH-Identity-Proof": {"proof"}})
			Expect(err).To(MatchError(core.ErrWorkloadDenied))
			_, err = signedCall(v, pattern, path, key, private, core.WorkloadClaims{Scopes: []string{scope}, SessionID: "forged"}, map[string][]string{"X-GH-Identity-Proof": {"proof"}})
			Expect(err).To(MatchError(core.ErrWorkloadDenied))
			actor := session.Principal{AccountID: "account", SessionID: "session", TenantID: "tenant", PersonGID: "person"}
			gomock.InOrder(
				delegated.EXPECT().VerifyAccess(gomock.Any(), "bad").Return(session.Principal{}, errors.New("bad proof")),
				delegated.EXPECT().VerifyAccess(gomock.Any(), "missing-session").Return(session.Principal{AccountID: "account"}, nil),
				delegated.EXPECT().VerifyAccess(gomock.Any(), "missing-account").Return(session.Principal{SessionID: "session"}, nil),
				delegated.EXPECT().VerifyAccess(gomock.Any(), "denied").Return(actor, nil),
				delegated.EXPECT().AuthorizeProtectedRoute(gomock.Any(), "GET /identity/protected/auth/me").Return(errors.New("denied")),
			)
			for _, proof := range []string{"bad", "missing-session", "missing-account", "denied"} {
				_, err := signedCall(v, pattern, path, key, private, core.WorkloadClaims{Scopes: []string{scope}}, map[string][]string{"X-GH-Identity-Proof": {proof}})
				Expect(err).To(MatchError(core.ErrWorkloadDenied))
			}
		})
	})
}
