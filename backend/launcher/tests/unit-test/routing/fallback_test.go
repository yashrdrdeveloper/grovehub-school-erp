package routing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/launcher/routing"
)

var _ = Describe("Fallback routing", func() {
	It("returns the standard not-found response for an unknown route", func() {
		recorder := serveRequest(newFallbackMux(), http.MethodGet, unknownPath)
		response := decodeRouteError(recorder)

		Expect(recorder.Code).To(Equal(http.StatusNotFound))
		Expect(response.Success).To(BeFalse())
		Expect(response.ErrorCode).To(Equal(coreerrors.Code(expectedRouteErrorCode)))
		Expect(response.ErrorMessage).To(Equal(expectedRouteErrorReason))
		Expect(response.HTTPErrorCode).To(Equal(http.StatusNotFound))
	})

	It("does not override a registered route", func() {
		mux := http.NewServeMux()
		mux.HandleFunc("GET "+registeredIdentityPath, func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNoContent)
		})
		routing.RegisterFallback(mux)

		recorder := serveRequest(mux, http.MethodGet, registeredIdentityPath)
		Expect(recorder.Code).To(Equal(http.StatusNoContent))
	})
})
