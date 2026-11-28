package helpers_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/launcher/internal/src/services/gateway/operations/external/helpers"
)

// @Testing suite=LauncherExternalHelpers
func ResponseSpecs() {
	Describe("Reject", func() {
		It("uses consistent public error envelopes without caching or revealing policy internals", func() {
			for _, tc := range []struct {
				status int
				code   string
			}{
				{http.StatusUnauthorized, "GATEWAY_UNAUTHENTICATED"},
				{http.StatusForbidden, "GATEWAY_FORBIDDEN"},
				{http.StatusTooManyRequests, "GATEWAY_RATE_LIMITED"},
			} {
				request := httptest.NewRequest(http.MethodGet, "/", nil)
				response := httptest.NewRecorder()
				helpers.Reject(response, request, tc.status)
				Expect(response.Code).To(Equal(tc.status))
				Expect(response.Header().Get("Cache-Control")).To(Equal("no-store"))
				Expect(response.Body.String()).To(ContainSubstring(tc.code))
			}

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			helpers.RejectForbidden(response, request, "private-policy-detail")
			Expect(response.Code).To(Equal(http.StatusForbidden))
			Expect(response.Body.String()).NotTo(ContainSubstring("private-policy-detail"))
		})
	})
}
