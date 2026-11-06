package controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	integrationconfig "grovehub/services/integration/config"
	integrationcontroller "grovehub/services/integration/controller/impl"
	integrationrequest "grovehub/services/integration/dto/request"
	integrationresponse "grovehub/services/integration/dto/response"
	"grovehub/services/integration/providers"
	"grovehub/services/integration/providers/mailprovider/providers/gmail"
	integrationroutes "grovehub/services/integration/routes"
)

type serviceStub struct{}

func (*serviceStub) Providers(context.Context) []providers.Descriptor {
	return []providers.Descriptor{{ID: "GMAIL"}}
}
func (*serviceStub) Provider(context.Context, string) (integrationresponse.ProviderDetails, error) {
	return integrationresponse.ProviderDetails{}, nil
}
func (*serviceStub) Create(context.Context, string, integrationrequest.CreateConnection) (integrationresponse.Connection, error) {
	return integrationresponse.Connection{}, nil
}
func (*serviceStub) Get(context.Context, string, string) (integrationresponse.Connection, error) {
	return integrationresponse.Connection{}, nil
}
func (*serviceStub) List(context.Context, string, string) ([]integrationresponse.Connection, error) {
	return nil, nil
}
func (*serviceStub) BeginOAuth(context.Context, string, string) (integrationresponse.OAuthStart, error) {
	return integrationresponse.OAuthStart{}, nil
}
func (*serviceStub) FinishOAuth(context.Context, string, string, string) (integrationresponse.OAuthCompletion, error) {
	return integrationresponse.OAuthCompletion{}, nil
}
func (*serviceStub) Watch(context.Context, string, string) (integrationresponse.Watch, error) {
	return integrationresponse.Watch{}, nil
}
func (*serviceStub) ProcessPush(context.Context, gmail.PushNotification) error { return nil }
func (*serviceStub) RenewExpiring(context.Context) error                       { return nil }
func (*serviceStub) RunTest(context.Context, string, string, string) (integrationresponse.ConnectionTest, error) {
	return integrationresponse.ConnectionTest{}, nil
}
func (*serviceStub) ListTests(context.Context, string, string) ([]integrationresponse.ConnectionTest, error) {
	return nil, nil
}
func (*serviceStub) ExecuteAction(context.Context, string, string, providers.Action, json.RawMessage) (any, error) {
	return nil, nil
}
func (*serviceStub) Disconnect(context.Context, string, string) (integrationresponse.Connection, error) {
	return integrationresponse.Connection{}, nil
}

var _ = Describe("Integration HTTP authorization", func() {
	It("allows public provider discovery", func() {
		c := integrationcontroller.NewConnectionController(&serviceStub{}, &serviceStub{}, integrationconfig.GmailPush{})
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/integrations/providers", nil)
		c.Providers(recorder, req)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body.String()).To(ContainSubstring("GMAIL"))
	})
	It("denies unauthorized OAuth starts and mail tests", func() {
		c := integrationcontroller.NewConnectionController(&serviceStub{}, &serviceStub{}, integrationconfig.GmailPush{})
		for _, testCase := range []struct {
			route   string
			handler http.HandlerFunc
		}{
			{"/integrations/connections/abc/oauth/start", c.BeginOAuth},
			{"/integrations/connections/abc/test", c.RunTest},
			{"/integrations/connections/abc/tests", c.ListTests},
			{"/integrations/connections/abc/actions/mail.send", c.ExecuteAction},
			{"/integrations/connections/abc/watch", c.Watch},
		} {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, testCase.route, nil)
			testCase.handler(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusForbidden))
		}
	})
	It("rejects unauthenticated Gmail push", func() {
		c := integrationcontroller.NewConnectionController(&serviceStub{}, &serviceStub{}, integrationconfig.GmailPush{Topic: "projects/test/topics/mail", Audience: "https://example.test/push", ServiceAccount: "push@example.test"})
		response := httptest.NewRecorder()
		c.GmailPush(response, httptest.NewRequest(http.MethodPost, "/integrations/webhooks/gmail", strings.NewReader("{}")))
		Expect(response.Code).To(Equal(http.StatusUnauthorized))
	})
	It("registers all connection URLs in the connections resource group", func() {
		c := integrationcontroller.NewConnectionController(&serviceStub{}, &serviceStub{}, integrationconfig.GmailPush{})
		mux := http.NewServeMux()
		integrationroutes.RegisterHTTP(mux, c)

		for _, route := range []struct {
			method, path string
			want         int
		}{
			{http.MethodGet, "/integrations/providers", http.StatusOK},
			{http.MethodGet, "/integrations/connections/oauth/gmail/callback", http.StatusOK},
			{http.MethodGet, "/integrations/connections/oauth/outlook/callback", http.StatusOK},
			{http.MethodPost, "/integrations/connections/example/oauth/start", http.StatusForbidden},
			{http.MethodPost, "/integrations/connections/example/actions/mail.send", http.StatusForbidden},
			{http.MethodPost, "/integrations/connections/example/test", http.StatusForbidden},
			{http.MethodPost, "/integrations/connections/example/watch", http.StatusForbidden},
			{http.MethodGet, "/integrations/oauth/gmail/callback", http.StatusNotFound},
		} {
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(route.method, route.path, nil))
			Expect(response.Code).To(Equal(route.want), route.path)
		}
	})
	It("permits a matching development actor on loopback only when explicitly enabled",func(){
  dev:=integrationconfig.DevAuth{
   Enabled:true,
   TenantID:"018f5d04-315c-75a9-b890-df0d02e787c1",
   ActorGID:"01a1b210-402c-7622-a33c-f82ceba26c7d",
  }
  c:=integrationcontroller.NewConnectionController(&serviceStub{},&serviceStub{},integrationconfig.GmailPush{},dev)
  requestBody:=`{"provider":"GMAIL","ownerGid":"01a1b210-402c-7622-a33c-f82ceba26c7d"}`

  local:=httptest.NewRequest(http.MethodPost,"/integrations/connections",strings.NewReader(requestBody))
  local.RemoteAddr="127.0.0.1:50231"
  response:=httptest.NewRecorder()
  c.Create(response,local)
  Expect(response.Code).To(Equal(http.StatusCreated))

  remote:=httptest.NewRequest(http.MethodPost,"/integrations/connections",strings.NewReader(requestBody))
  remote.RemoteAddr="198.51.100.10:50231"
  remote.Header.Set("X-Forwarded-For","127.0.0.1")
  response=httptest.NewRecorder()
  c.Create(response,remote)
  Expect(response.Code).To(Equal(http.StatusForbidden))

  forged:=httptest.NewRequest(http.MethodPost,"/integrations/connections",strings.NewReader(`{"provider":"GMAIL","ownerGid":"018f5d04-315c-75a9-b890-df0d02e787c1"}`))
  forged.RemoteAddr="[::1]:50231"
  response=httptest.NewRecorder()
  c.Create(response,forged)
  Expect(response.Code).To(Equal(http.StatusForbidden))
 })
 It("denies connection creation without a verified principal", func() {
		c := integrationcontroller.NewConnectionController(&serviceStub{}, &serviceStub{}, integrationconfig.GmailPush{})
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/integrations/connections", strings.NewReader(`{"provider":"GMAIL","ownerGid":"a"}`))
		req.Header.Set("X-Tenant-ID", "spoofed")
		c.Create(recorder, req)
		Expect(recorder.Code).To(Equal(http.StatusForbidden))
	})
})
