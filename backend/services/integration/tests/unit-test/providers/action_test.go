package providers_test

import (
	"context"
	"encoding/json"
 "encoding/base64"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	integrationconfig "grovehub/services/integration/config"
	"grovehub/services/integration/providers"
	providerimpl "grovehub/services/integration/providers/impl"
	"grovehub/services/integration/providers/mailprovider/providers/gmail"
	"grovehub/services/integration/providers/mailprovider/providers/outlook"
)

type actionFake struct{}

func (actionFake) Execute(_ context.Context, _ providers.Session, _ json.RawMessage) (any, error) {
	return map[string]bool{"executed": true}, nil
}

var _ = Describe("Generic action strategy registry", func() {
	It("routes only registered capabilities without provider-name branching", func() {
		registry := providerimpl.NewActionRegistry()
		Expect(registry.Register("CUSTOM", providers.Action("lead.import"), actionFake{})).To(Succeed())
		result, err := registry.Execute(context.Background(), providers.Session{ProviderID: "CUSTOM"}, providers.Action("lead.import"), json.RawMessage("{}"))
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(map[string]bool{"executed": true}))
		_, err = registry.Execute(context.Background(), providers.Session{ProviderID: "CUSTOM"}, providers.ActionMailSend, nil)
		Expect(err).To(MatchError(providers.ErrUnsupportedAction))
	})
	It("rejects duplicate action registrations", func() {
		registry := providerimpl.NewActionRegistry()
		Expect(registry.Register("GMAIL", providers.ActionMailSend, actionFake{})).To(Succeed())
		Expect(registry.Register("GMAIL", providers.ActionMailSend, actionFake{})).To(HaveOccurred())
	})
	It("normalizes rate limits as retryable errors", func() {
		err := providers.NormalizeStatus(429)
		Expect(providers.HTTPStatus(err)).To(Equal(429))
		Expect(err.Error()).NotTo(ContainSubstring("token"))
	})
	It("rejects malicious Outlook pagination URLs without remote requests", func() {
		client := outlook.NewClient(integrationconfig.Outlook{})
		_, err := client.Messages(context.Background(), "test", "https://evil.example/v1.0/me/messages?$skip=1")
		Expect(err).To(MatchError(providers.ErrInvalidActionInput))
		_, err = client.Delta(context.Background(), "test", "https://graph.microsoft.com.evil.example/v1.0/me/mailFolders/inbox/messages/delta")
		Expect(err).To(MatchError(providers.ErrInvalidActionInput))
	})
	It("rejects obsolete raw Gmail send payloads without contacting Google", func() {
  registry:=providerimpl.NewActionRegistry()
  client:=gmail.NewClient(integrationconfig.Gmail{})
  Expect(gmail.RegisterActions(registry,client)).To(Succeed())
  session:=providers.Session{ProviderID:"GMAIL",AccessToken:"dummy"}
  for _,raw:=range []string{"base64url-encoded-RFC5322","%%%","",base64.RawURLEncoding.EncodeToString([]byte("hello"))}{
   body,_:=json.Marshal(map[string]string{"raw":raw})
   _,err:=registry.Execute(context.Background(),session,providers.ActionMailSend,body)
   Expect(err).To(MatchError(providers.ErrInvalidActionInput))
  }
 })
	It("rejects malformed Gmail raw mail without remote requests", func() {
		client := gmail.NewClient(integrationconfig.Gmail{})
		_, err := client.SendRaw(context.Background(), gmail.Token{AccessToken: "dummy"}, "%%%")
		Expect(err).To(HaveOccurred())
	})
})
