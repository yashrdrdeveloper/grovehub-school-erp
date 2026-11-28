package module_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coremodule "grovehub/libs/core/module"
)

// Service-owned implementations can extend the shared contract.
type serviceModule interface {
	coremodule.Module
	Name() string
}

type exampleModule struct {
	closed bool
}

var _ serviceModule = (*exampleModule)(nil)

func (m *exampleModule) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /example", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func (m *exampleModule) Close()       { m.closed = true }
func (m *exampleModule) Name() string { return "example" }

var _ = Describe("Module facade", func() {
	It("allows service implementations to register routes and close via the shared interface", func() {
		service := &exampleModule{}
		var lifecycle coremodule.Module = service
		mux := http.NewServeMux()
		lifecycle.RegisterHTTP(mux)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/example", nil))
		Expect(recorder.Code).To(Equal(http.StatusNoContent))
		Expect(service.closed).To(BeFalse())

		lifecycle.Close()
		Expect(service.closed).To(BeTrue())
		Expect(service.Name()).To(Equal("example"))
	})
})
