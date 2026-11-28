package gateway_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	external "grovehub/launcher/internal/src/services/gateway/operations/external/impl"
)

// @Testing suite=LauncherExternalGateway
func TierSpecs() {
	Describe("Unknown tier", func() {
		It("fails closed before dispatching a registered endpoint outside the three supported tiers", func() {
			mux := http.NewServeMux()
			mux.HandleFunc("GET /foundation/unknown/persons", func(w http.ResponseWriter, _ *http.Request) {
				Fail("unclassified handler invoked")
			})
			w := httptest.NewRecorder()
			external.New(mux, nil, nil, "").ServeHTTP(w, httptest.NewRequest("GET", "/foundation/unknown/persons", nil))
			Expect(w.Code).To(Equal(http.StatusForbidden))
		})
	})
}
