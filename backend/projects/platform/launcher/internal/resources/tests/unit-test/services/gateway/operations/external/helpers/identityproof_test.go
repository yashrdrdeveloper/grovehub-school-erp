package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/launcher/internal/resources/tests/helpers/workload"
	"grovehub/launcher/internal/src/services/gateway/operations/external/helpers"
	"grovehub/libs/platform/session"
)

// @Testing suite=LauncherExternalHelpers
func IdentityProofSpecs() {
	Describe("AdaptIdentityProof", func() {
		It("maps every verified School identity-proof route to its credential format", func() {
			ctx, err := fixture.VerifiedContext("school")
			Expect(err).NotTo(HaveOccurred())
			for _, tc := range []struct{ pattern, cookie, prefix string }{
				{"POST /identity/restricted/auth/refresh", session.RefreshCookie, ""},
				{"GET /identity/restricted/auth/recovery/validate", "", "Recovery"},
				{"GET /identity/restricted/auth/invitations/validate", "", "Invitation"},
				{"POST /identity/restricted/auth/recovery/exchange", "", "Recovery"},
				{"POST /identity/restricted/auth/invitations/exchange", "", "Invitation"},
				{"GET /identity/restricted/auth/reset-password/status", session.ResetCookie, ""},
				{"POST /identity/restricted/auth/reset-password", session.ResetCookie, ""},
				{"GET /identity/restricted/accounts/setup/status", session.SetupCookie, ""},
				{"POST /identity/restricted/accounts/setup", session.SetupCookie, ""},
			} {
				method, path, _ := strings.Cut(tc.pattern, " ")
				request := httptest.NewRequest(method, path, nil).WithContext(ctx)
				request.Pattern = tc.pattern
				request.Header.Set("X-GH-Identity-Proof", "valid-proof")
				request.Header.Set("Authorization", "Bearer workload-token")
				request.AddCookie(&http.Cookie{Name: "untrusted", Value: "value"})
				Expect(helpers.AdaptIdentityProof(request)).To(BeTrue(), tc.pattern)
				Expect(request.Header.Get("X-GH-Identity-Proof")).To(BeEmpty())
				_, err := request.Cookie("untrusted")
				Expect(err).To(MatchError(http.ErrNoCookie))
				if tc.cookie != "" {
					cookie, err := request.Cookie(tc.cookie)
					Expect(err).NotTo(HaveOccurred())
					Expect(cookie.Value).To(Equal("valid-proof"))
					Expect(request.Header.Get("Authorization")).To(BeEmpty())
				} else {
					Expect(request.Header.Get("Cookie")).To(BeEmpty())
					Expect(request.Header.Get("Authorization")).To(Equal(tc.prefix + " valid-proof"))
				}
			}
		})
		It("rejects unverified and malformed proofs and leaves unrelated routes unchanged", func() {
			request := httptest.NewRequest(http.MethodPost, "/identity/restricted/auth/refresh", nil)
			request.Pattern = "POST /identity/restricted/auth/refresh"
			request.Header.Set("X-GH-Identity-Proof", "valid-proof")
			Expect(helpers.AdaptIdentityProof(request)).To(BeFalse())
			foreign, err := fixture.VerifiedContext("internal-tool")
			Expect(err).NotTo(HaveOccurred())
			Expect(helpers.AdaptIdentityProof(request.WithContext(foreign))).To(BeFalse())
			school, err := fixture.VerifiedContext("school")
			Expect(err).NotTo(HaveOccurred())
			for _, values := range [][]string{nil, {""}, {"a", "b"}, {strings.Repeat("x", 2049)}, {"a b"}, {"a,b"}, {"a;b"}} {
				r := request.Clone(school)
				r.Header.Del("X-GH-Identity-Proof")
				for _, value := range values {
					r.Header.Add("X-GH-Identity-Proof", value)
				}

				Expect(helpers.AdaptIdentityProof(r)).To(BeFalse())
			}

			request.Pattern = "POST /identity/restricted/auth/login"
			Expect(helpers.AdaptIdentityProof(request)).To(BeTrue())
			Expect(request.Header.Get("X-GH-Identity-Proof")).To(Equal("valid-proof"))
		})
	})
}
