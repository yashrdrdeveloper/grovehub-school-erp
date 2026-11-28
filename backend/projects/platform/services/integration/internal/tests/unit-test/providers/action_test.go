package providers_test

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    providerhelpers "grovehub/services/integration/internal/app/providers/helpers"
    mailconstants "grovehub/services/integration/internal/app/providers/mailprovider/constants"
    execute "grovehub/services/integration/internal/app/providers/operations/execute"
	"context"
	"encoding/json"
 "encoding/base64"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	integrationconfig "grovehub/services/integration/internal/app/config"
	"grovehub/services/integration/internal/app/externalService/mail/gmail"
	"grovehub/services/integration/internal/app/externalService/mail/outlook"
)

type actionFake struct{}

func (actionFake) Execute(_ context.Context, _ providerdto.Session, _ json.RawMessage) (any, error) {
	return map[string]bool{"executed": true}, nil
}

var _ = Describe("Generic action strategy registry", func() {
	It("routes only registered capabilities without provider-name branching", func() {
		registry := execute.NewActionRegistry()
		Expect(registry.Register("CUSTOM", providerdto.Action("lead.import"), actionFake{})).To(Succeed())
		result, err := registry.Execute(context.Background(), providerdto.Session{ProviderID: "CUSTOM"}, providerdto.Action("lead.import"), json.RawMessage("{}"))
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(map[string]bool{"executed": true}))
		_, err = registry.Execute(context.Background(), providerdto.Session{ProviderID: "CUSTOM"}, mailconstants.MailSend, nil)
		Expect(err).To(MatchError(providerhelpers.ErrUnsupportedAction))
	})
	It("rejects duplicate action registrations", func() {
		registry := execute.NewActionRegistry()
		Expect(registry.Register("GMAIL", mailconstants.MailSend, actionFake{})).To(Succeed())
		Expect(registry.Register("GMAIL", mailconstants.MailSend, actionFake{})).To(HaveOccurred())
	})
	It("normalizes rate limits as retryable errors", func() {
		err := providerhelpers.NormalizeStatus(429)
		Expect(providerhelpers.HTTPStatus(err)).To(Equal(429))
		Expect(err.Error()).NotTo(ContainSubstring("token"))
	})
	It("rejects malicious Outlook pagination URLs without remote requests", func() {
		client := outlook.NewClient(integrationconfig.Outlook{})
		_, err := client.Messages(context.Background(), "test", "https://evil.example/v1.0/me/messages?$skip=1")
		Expect(err).To(MatchError(providerhelpers.ErrInvalidActionInput))
		_, err = client.Delta(context.Background(), "test", "https://graph.microsoft.com.evil.example/v1.0/me/mailFolders/inbox/messages/delta")
		Expect(err).To(MatchError(providerhelpers.ErrInvalidActionInput))
	})
	It("rejects obsolete raw Gmail send payloads without contacting Google", func() {
  registry:=execute.NewActionRegistry()
  client:=gmail.NewClient(integrationconfig.Gmail{})
  Expect(gmail.RegisterActions(registry,client)).To(Succeed())
  session:=providerdto.Session{ProviderID:"GMAIL",AccessToken:"dummy"}
  for _,raw:=range []string{"base64url-encoded-RFC5322","%%%","",base64.RawURLEncoding.EncodeToString([]byte("hello"))}{
   body,_:=json.Marshal(map[string]string{"raw":raw})
   _,err:=registry.Execute(context.Background(),session,mailconstants.MailSend,body)
   Expect(err).To(MatchError(providerhelpers.ErrInvalidActionInput))
  }
 })
	It("rejects malformed Gmail raw mail without remote requests", func() {
		client := gmail.NewClient(integrationconfig.Gmail{})
		_, err := client.SendRaw(context.Background(), gmail.Token{AccessToken: "dummy"}, "%%%")
		Expect(err).To(HaveOccurred())
	})
})
