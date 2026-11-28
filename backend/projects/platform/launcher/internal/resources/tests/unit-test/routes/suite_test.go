package routing_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/gomega"

	corehttp "grovehub/libs/core/http"
	routing "grovehub/launcher/internal/src/routes"
)

const (
	unknownPath              = "/abc"
	registeredIdentityPath   = "/identity/health"
	expectedRouteErrorCode   = "PLATFORM_ROUTE_NOT_FOUND"
	expectedRouteErrorReason = "API route not found"
)

func newFallbackMux() *http.ServeMux {
	mux := http.NewServeMux()
	routing.RegisterFallback(mux)
	return mux
}

func serveRequest(handler http.Handler, method string, path string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func decodeRouteError(recorder *httptest.ResponseRecorder) corehttp.Error[struct{}] {
	var response corehttp.Error[struct{}]
	Expect(json.NewDecoder(recorder.Body).Decode(&response)).To(Succeed())
	return response
}
