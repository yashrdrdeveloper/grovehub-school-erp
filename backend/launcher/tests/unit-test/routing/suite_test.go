package routing_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	responsedto "grovehub/libs/core/dto/response"
	"grovehub/launcher/routing"
)

const (
	unknownPath              = "/abc"
	registeredIdentityPath   = "/identity/health"
	expectedRouteErrorCode   = "PLATFORM_ROUTE_NOT_FOUND"
	expectedRouteErrorReason = "API route not found"
)

func TestLauncherRoutingUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Launcher Routing Unit Tests")
}

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

func decodeRouteError(recorder *httptest.ResponseRecorder) responsedto.Error[struct{}] {
	var response responsedto.Error[struct{}]
	Expect(json.NewDecoder(recorder.Body).Decode(&response)).To(Succeed())
	return response
}
