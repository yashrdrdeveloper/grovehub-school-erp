package providers_test

import (
 "net/url"

 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"

 integrationconfig "grovehub/services/integration/config"
 "grovehub/services/integration/providers/mailprovider/providers/gmail"
)
var _=Describe("Gmail OAuth authorization URL",func(){
 It("binds a one-time state, the registered callback, and required API scopes",func(){
  client:=gmail.NewClient(integrationconfig.Gmail{
   ClientID:"test-client",ClientSecret:"test-secret",
   RedirectURI:"https://groovehub.example/integrations/oauth/gmail/callback",
  })
  parsed,err:=url.Parse(client.AuthorizationURL("random-secret-state"))
  Expect(err).NotTo(HaveOccurred())
  Expect(parsed.Host).To(Equal("accounts.google.com"))
  Expect(parsed.Query().Get("state")).To(Equal("random-secret-state"))
  Expect(parsed.Query().Get("response_type")).To(Equal("code"))
  Expect(parsed.Query().Get("access_type")).To(Equal("offline"))
  Expect(parsed.Query().Get("redirect_uri")).To(Equal("https://groovehub.example/integrations/oauth/gmail/callback"))
  Expect(parsed.Query().Get("scope")).To(ContainSubstring("gmail.readonly"))
  Expect(parsed.Query().Get("scope")).To(ContainSubstring("gmail.send"))
 })
})
