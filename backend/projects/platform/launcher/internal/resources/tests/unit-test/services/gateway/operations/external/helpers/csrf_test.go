package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/internal/session"
	"grovehub/launcher/internal/src/services/gateway/operations/external/helpers"
)

// @Testing suite=LauncherExternalHelpers
func CSRFSpecs() {
	Describe("RequiresCSRF", func() {
		It("requires matching browser CSRF credentials only on unsafe session routes", func() {
			for _, tc := range []struct {
				method, path string
				required bool
			}{
				{"GET", "/identity/protected/auth/me", false},
				{"HEAD", "/identity/protected/auth/me", false},
				{"POST", "/identity/protected/auth/logout", true},
				{"POST", "/identity/public/auth/refresh", true},
				{"POST", "/identity/public/auth/reset-password", true},
				{"POST", "/identity/public/accounts/setup", true},
				{"POST", "/identity/public/auth/login", false},
			} {
				request := httptest.NewRequest(tc.method, tc.path, nil)
				request.Pattern = tc.method+" "+tc.path
				Expect(helpers.RequiresCSRF(request)).To(Equal(tc.required), tc.path)
			}
		})
	})
	Describe("ValidCSRF", func() {
		It("rejects missing, empty, mismatched and oversized credentials", func() {
			for _, tc := range []struct{
				cookie, header string
				hasCookie, allowed bool
			}{
				{"", "token", false, false},
				{"token", "", true, false},
				{"", "token", true, false},
				{"token", "other", true, false},
				{strings.Repeat("a",257), strings.Repeat("a",257), true, false},
				{"token", "token", true, true},
			} {
				request := httptest.NewRequest(http.MethodPost, "/", nil)
				if tc.hasCookie { request.AddCookie(&http.Cookie{Name:session.CSRFName, Value:tc.cookie}) }
				request.Header.Set("X-CSRF-Token", tc.header)
				Expect(helpers.ValidCSRF(request)).To(Equal(tc.allowed))
			}
		})
	})
}
