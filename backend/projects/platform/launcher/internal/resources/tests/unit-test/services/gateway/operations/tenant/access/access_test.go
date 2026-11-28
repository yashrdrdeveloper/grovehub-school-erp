package access_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/launcher/internal/src/services/gateway/operations/tenant/access"
	core "grovehub/libs/core/gateway"
	"grovehub/libs/platform/session"
)

func tenantConfig(origin string) access.TenantAccessConfig {
	_, private, err := ed25519.GenerateKey(rand.Reader)
	Expect(err).NotTo(HaveOccurred())
	return access.TenantAccessConfig{Origin: origin, KeyID: "tenant", PrivateKey: base64.RawURLEncoding.EncodeToString(private)}
}

// @Testing suite=LauncherTenantAccess
func TenantAccessSpecs() {
	Describe("Tenant access boundary", func() {
		It("rejects invalid origins and credentials and preserves optional unconfigured access", func() {
			check, err := access.NewTenantAccess(access.TenantAccessConfig{})
			Expect(err).NotTo(HaveOccurred())
			Expect(check).To(BeNil())
			for _, origin := range []string{"http://remote.example", "https://example/path", "https://example?query=1", "https://user@example", "https://example#fragment", "%"} {
				_, err := access.NewTenantAccess(tenantConfig(origin))
				Expect(err).To(MatchError("invalid tenant access origin"))
			}

			for _, config := range []access.TenantAccessConfig{
				{Origin: "https://school.example", KeyID: "key", PrivateKey: "not-base64"},
				{Origin: "https://school.example", KeyID: "key", PrivateKey: "AA"},
			} {
				_, err := access.NewTenantState(config)
				Expect(err).To(MatchError("invalid tenant access credential"))
			}
		})
		It("restricts outbound credentials to exact body-free canonical registry requests", func() {
			path := "/school/restricted/registry/schools/" + uuid.NewString() + "/access-check"
			r := httptest.NewRequest("GET", path, nil)
			claims, err := access.TenantAccessPolicy(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(claims.Scopes).To(Equal([]string{access.TenantAccessScope}))
			_, err = access.TenantAccessPolicy(nil)
			Expect(err).To(MatchError(core.ErrWorkloadDenied))
			for _, mutate := range []func(*http.Request){
				func(r *http.Request) {
					r.URL = nil
				},
				func(r *http.Request) {
					r.Method = "POST"
				},
				func(r *http.Request) {
					r.URL.Path = "/wrong"
				},
				func(r *http.Request) {
					r.URL.RawQuery = "query=1"
				},
				func(r *http.Request) {
					r.URL.RawPath = "/encoded"
				},
				func(r *http.Request) {
					r.URL.ForceQuery = true
				},
				func(r *http.Request) {
					r.ContentLength = 1
				},
			} {
				request := httptest.NewRequest("GET", path, nil)
				mutate(request)
				_, err := access.TenantAccessPolicy(request)
				Expect(err).To(MatchError(core.ErrWorkloadDenied))
			}
		})
		It("fails closed when client construction, request construction or network transport fails", func() {
			config := tenantConfig("https://school.example")
			failure := errors.New("client unavailable")
			_, err := access.NewTenantStateWithDependencies(config, func(*core.WorkloadTransport) (*http.Client, error) {
				return nil, failure
			}, http.NewRequestWithContext)
			Expect(err).To(MatchError(failure))
			lookup, err := access.NewTenantStateWithDependencies(config, core.NewWorkloadHTTPClient, func(context.Context, string, string, io.Reader) (*http.Request, error) {
				return nil, errors.New("request construction failed")
			})
			Expect(err).NotTo(HaveOccurred())
			_, err = lookup(context.Background(), uuid.NewString())
			Expect(err).To(MatchError(session.ErrTenantDenied))
			server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
			config = tenantConfig(server.URL)
			lookup, err = access.NewTenantState(config)
			Expect(err).NotTo(HaveOccurred())
			server.Close()
			_, err = lookup(context.Background(), uuid.NewString())
			Expect(err).To(MatchError(session.ErrTenantDenied))
		})
		It("rejects response bodies and ambiguous generation headers", func() {
			for _, tc := range []struct {
				body        string
				generations []string
			}{
				{"unexpected", []string{"1"}},
				{"", []string{"1", "2"}},
			} {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					for _, value := range tc.generations {
						w.Header().Add("X-GH-Session-Generation", value)
					}

					w.WriteHeader(http.StatusForbidden)
					_, _ = io.Copy(w, strings.NewReader(tc.body))
				}))
				lookup, err := access.NewTenantState(tenantConfig(server.URL))
				Expect(err).NotTo(HaveOccurred())
				_, err = lookup(context.Background(), uuid.NewString())
				server.Close()
				Expect(err).To(MatchError(session.ErrTenantDenied))
			}
		})
	})
}
