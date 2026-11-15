package controller_test

import (
    providerdto "grovehub/services/integration/internal/app/dto/provider"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	integrationconfig "grovehub/services/integration/config"
	integrationcontroller "grovehub/services/integration/internal/app/controller/connection/impl"
    combinedcontroller "grovehub/services/integration/internal/app/controller/impl"
    providercontroller "grovehub/services/integration/internal/app/controller/provider/impl"
	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	
	"grovehub/services/integration/internal/app/providers"
	"grovehub/services/integration/internal/app/externalService/mail/gmail"
	integrationroutes "grovehub/services/integration/internal/app/routes"
)

type serviceStub struct{}

func (*serviceStub) Providers(context.Context) []providers.Descriptor {
	return []providers.Descriptor{{ID: "GMAIL"}}
}
func (*serviceStub) Provider(context.Context, string) (providerdto.ProviderDetails, error) {
	return providerdto.ProviderDetails{}, nil
}
func (*serviceStub) Create(context.Context, string, connectiondto.CreateConnection) (connectiondto.Connection, error) {
	return connectiondto.Connection{}, nil
}
func (*serviceStub) Get(context.Context, string, string) (connectiondto.Connection, error) {
	return connectiondto.Connection{}, nil
}
func (*serviceStub) List(context.Context, string, string) ([]connectiondto.Connection, error) {
	return nil, nil
}
func (*serviceStub) BeginOAuth(context.Context, string, string) (connectiondto.OAuthStart, error) {
	return connectiondto.OAuthStart{}, nil
}
func (*serviceStub) FinishOAuth(context.Context, string, string, string) (connectiondto.OAuthCompletion, error) {
	return connectiondto.OAuthCompletion{}, nil
}
func (*serviceStub) Watch(context.Context, string, string) (connectiondto.Watch, error) {
	return connectiondto.Watch{}, nil
}
func (*serviceStub) ProcessPush(context.Context, gmail.PushNotification) error { return nil }
func (*serviceStub) RenewExpiring(context.Context) error                       { return nil }
func (*serviceStub) RunTest(context.Context, string, string, string) (connectiondto.ConnectionTest, error) {
	return connectiondto.ConnectionTest{}, nil
}
func (*serviceStub) ListTests(context.Context, string, string) ([]connectiondto.ConnectionTest, error) {
	return nil, nil
}
func (*serviceStub) ExecuteAction(context.Context, string, string, providers.Action, json.RawMessage) (any, error) {
	return nil, nil
}
func (*serviceStub) Disconnect(context.Context, string, string) (connectiondto.Connection, error) {
	return connectiondto.Connection{}, nil
}

var _ = Describe("Integration HTTP authorization", func() {
	It("allows public provider discovery", func() {
		c := integrationcontroller.NewConnectionController(&serviceStub{}, integrationconfig.GmailPush{})
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/integrations/providers", nil)
		providercontroller.NewProviderController(&serviceStub{}).Providers(recorder, req)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body.String()).To(ContainSubstring("GMAIL"))
	})
	It("denies unauthorized OAuth starts and mail tests", func() {
		c := integrationcontroller.NewConnectionController(&serviceStub{}, integrationconfig.GmailPush{})
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
		c := integrationcontroller.NewConnectionController(&serviceStub{}, integrationconfig.GmailPush{Topic: "projects/test/topics/mail", Audience: "https://example.test/push", ServiceAccount: "push@example.test"})
		response := httptest.NewRecorder()
		c.GmailPush(response, httptest.NewRequest(http.MethodPost, "/integrations/webhooks/gmail", strings.NewReader("{}")))
		Expect(response.Code).To(Equal(http.StatusUnauthorized))
	})
	It("registers all connection URLs in the connections resource group", func() {
		c := integrationcontroller.NewConnectionController(&serviceStub{}, integrationconfig.GmailPush{})
		mux := http.NewServeMux()
		integrationroutes.RegisterHTTP(mux, combinedcontroller.New(c,providercontroller.NewProviderController(&serviceStub{})))

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
  c:=integrationcontroller.NewConnectionController(&serviceStub{},integrationconfig.GmailPush{},dev)
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
		c := integrationcontroller.NewConnectionController(&serviceStub{}, integrationconfig.GmailPush{})
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/integrations/connections", strings.NewReader(`{"provider":"GMAIL","ownerGid":"a"}`))
		req.Header.Set("X-Tenant-ID", "spoofed")
		c.Create(recorder, req)
		Expect(recorder.Code).To(Equal(http.StatusForbidden))
	})
})
