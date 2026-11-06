package providers_test

import (
 "net/url"

 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"

 integrationconfig "grovehub/services/integration/config"
 "grovehub/services/integration/providers/outlook"
)

var _=Describe("Outlook OAuth",func(){
 It("requests delegated Mail.Read, Mail.Send and offline_access with state",func(){
  client:=outlook.NewClient(integrationconfig.Outlook{
   Tenant:"common",ClientID:"test-client",
   ClientSecret:"test-secret",
   RedirectURI:"https://test.example/integrations/oauth/outlook/callback",
  })
  parsed,err:=url.Parse(client.AuthorizationURL("unique-state"))
  Expect(err).NotTo(HaveOccurred())
  Expect(parsed.Host).To(Equal("login.microsoftonline.com"))
  Expect(parsed.Query().Get("state")).To(Equal("unique-state"))
  Expect(parsed.Query().Get("scope")).To(ContainSubstring("offline_access"))
  Expect(parsed.Query().Get("scope")).To(ContainSubstring("Mail.Send"))
  Expect(parsed.Query().Get("redirect_uri")).To(Equal("https://test.example/integrations/oauth/outlook/callback"))
 })
})
