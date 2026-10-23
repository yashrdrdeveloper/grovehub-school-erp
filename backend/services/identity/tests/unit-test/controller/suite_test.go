package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	responsedto "grovehub/libs/core/dto/response"
	identityresponse "grovehub/services/identity/dto/v1/response"
)

const (
	healthPath            = "/identity/health"
	expectedHealthService = "identity"
	expectedHealthStatus  = "ok"
)

func TestIdentityControllerUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Identity Controller Unit Tests")
}

func newHealthRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, healthPath, nil)
}

func decodeHealthSuccess(recorder *httptest.ResponseRecorder) responsedto.Success[identityresponse.HealthData] {
	var payload responsedto.Success[identityresponse.HealthData]
	Expect(json.Unmarshal(recorder.Body.Bytes(), &payload)).To(Succeed())
	return payload
}
