package controller_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/controller"
)

var _ = Describe("Health controller", func() {
	It("returns the standard success envelope", func() {
		recorder := httptest.NewRecorder()
		controller.NewHealthController().ServeHTTP(recorder, newHealthRequest())
		payload := decodeHealthSuccess(recorder)

		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(payload.Success).To(BeTrue())
		Expect(payload.Data.Service).To(Equal(expectedHealthService))
		Expect(payload.Data.Status).To(Equal(expectedHealthStatus))
	})
})
