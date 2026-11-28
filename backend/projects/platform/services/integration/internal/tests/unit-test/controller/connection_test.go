package controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	clientprincipal "grovehub/services/integration/client/principal"
	integrationconfig "grovehub/services/integration/internal/app/config"
	publicconnections "grovehub/services/integration/internal/app/controller/connection/public/impl"
	protectedconnections "grovehub/services/integration/internal/app/controller/connection/protected/impl"
	providerscontroller "grovehub/services/integration/internal/app/controller/provider/public/impl"
	restrictedwebhooks "grovehub/services/integration/internal/app/controller/webhook/restricted/impl"
	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	providerresponse "grovehub/services/integration/internal/app/dto/provider"
	"grovehub/services/integration/internal/app/externalService/mail/gmail"
	providerdescriptor "grovehub/services/integration/internal/app/providers/dto"
	"grovehub/services/integration/internal/app/routes"
)

type serviceStub struct{}

func (*serviceStub) Providers(context.Context) []providerdescriptor.Descriptor {
	return []providerdescriptor.Descriptor{{ID: "GMAIL"}}
}

func (*serviceStub) Provider(context.Context, string) (providerresponse.ProviderDetails, error) {
	return providerresponse.ProviderDetails{}, nil
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

func (*serviceStub) ProcessPush(context.Context, gmail.PushNotification) error {
	return nil
}

func (*serviceStub) RenewExpiring(context.Context) error {
	return nil
}

func (*serviceStub) RunTest(context.Context, string, string, string) (connectiondto.ConnectionTest, error) {
	return connectiondto.ConnectionTest{}, nil
}

func (*serviceStub) ListTests(context.Context, string, string) ([]connectiondto.ConnectionTest, error) {
	return nil, nil
}

func (*serviceStub) ExecuteAction(context.Context, string, string, providerdescriptor.Action, json.RawMessage) (any, error) {
	return nil, nil
}

func (*serviceStub) ExecuteSystemMail(context.Context, string, json.RawMessage) (any, error) {
	return nil, nil
}

func (*serviceStub) ConfigureSystemMailSender(context.Context, string, string) error {
	return nil
}

func (*serviceStub) Disconnect(context.Context, string, string) (connectiondto.Connection, error) {
	return connectiondto.Connection{}, nil
}

var _ = Describe("Integration access-tier controllers", func() {
	var service *serviceStub

	BeforeEach(func() {
		service = &serviceStub{}
	})

	It("serves public provider discovery independently of connection ownership", func() {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/integrations/restricted/providers", nil)

		providerscontroller.NewProviderController(service).Providers(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body.String()).To(ContainSubstring("GMAIL"))
	})

	It("denies protected operations without a gateway-verified principal", func() {
		controller := protectedconnections.New(service)
		for _, testCase := range []struct {
			path string
			handler http.HandlerFunc
		}{
			{"/integrations/restricted/connections/abc/oauth/start", controller.BeginOAuth},
			{"/integrations/restricted/connections/abc/test", controller.RunTest},
			{"/integrations/restricted/connections/abc/tests", controller.ListTests},
			{"/integrations/restricted/connections/abc/actions/mail.send", controller.ExecuteAction},
			{"/integrations/restricted/connections/abc/watch", controller.Watch},
		} {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, testCase.path, nil)
			request.Header.Set("X-Tenant-ID", "forged")

			testCase.handler(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusForbidden), testCase.path)
		}
	})

	It("rejects a Gmail webhook without valid signed workload authorization", func() {
		configuration := integrationconfig.GmailPush{
			Topic: "projects/test/topics/mail",
			Audience: "https://example.test/push",
			ServiceAccount: "push@example.test",
		}
		controller := restrictedwebhooks.New(service, configuration)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/integrations/restricted/webhooks/gmail", strings.NewReader("{}"))

		controller.GmailPush(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("registers only restricted Platform endpoints", func() {
		mux := http.NewServeMux()
		routes.RegisterHTTP(
			mux,
			publicconnections.New(service),
			protectedconnections.New(service),
			providerscontroller.NewProviderController(service),
			restrictedwebhooks.New(service, integrationconfig.GmailPush{}),
		)
		for _, testCase := range []struct {
			method string
			path string
			status int
		}{
			{http.MethodGet, "/integrations/restricted/providers", http.StatusOK},
			{http.MethodGet, "/integrations/restricted/connections/oauth/gmail/callback", http.StatusOK},
			{http.MethodGet, "/integrations/restricted/connections/oauth/outlook/callback", http.StatusOK},
			{http.MethodPost, "/integrations/restricted/connections/example/oauth/start", http.StatusForbidden},
			{http.MethodPost, "/integrations/restricted/connections/example/actions/mail.send", http.StatusForbidden},
			{http.MethodPost, "/integrations/restricted/connections/example/test", http.StatusForbidden},
			{http.MethodPost, "/integrations/restricted/connections/example/watch", http.StatusForbidden},
			{http.MethodGet, "/integrations/connections/oauth/gmail/callback", http.StatusNotFound},
			{http.MethodGet, "/integrations/public/providers", http.StatusNotFound},
			{http.MethodPost, "/integrations/protected/connections/example/test", http.StatusNotFound},
		} {
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(testCase.method, testCase.path, nil))
			Expect(recorder.Code).To(Equal(testCase.status), testCase.path)
		}
	})

	It("uses the verified caller context and never trusts a forged tenant header", func() {
		controller := protectedconnections.New(service)
		requestBody := `{"provider":"GMAIL","ownerGid":"01a1b210-402c-7622-a33c-f82ceba26c7d"}`

		verified := httptest.NewRequest(
			http.MethodPost,
			"/integrations/restricted/connections",
			strings.NewReader(requestBody),
		)
		verified = verified.WithContext(clientprincipal.WithVerifiedPrincipal(
			verified.Context(),
			clientprincipal.VerifiedPrincipal{
				TenantID: "018f5d04-315c-75a9-b890-df0d02e787c1",
				ActorGID: "01a1b210-402c-7622-a33c-f82ceba26c7d",
			},
		))
		recorder := httptest.NewRecorder()
		controller.Create(recorder, verified)
		Expect(recorder.Code).To(Equal(http.StatusCreated))

		forged := httptest.NewRequest(
			http.MethodPost,
			"/integrations/restricted/connections",
			strings.NewReader(requestBody),
		)
		forged.Header.Set("X-Tenant-ID", "018f5d04-315c-75a9-b890-df0d02e787c1")
		forged.Header.Set("X-Forwarded-For", "127.0.0.1")
		recorder = httptest.NewRecorder()
		controller.Create(recorder, forged)
		Expect(recorder.Code).To(Equal(http.StatusForbidden))
	})
})
