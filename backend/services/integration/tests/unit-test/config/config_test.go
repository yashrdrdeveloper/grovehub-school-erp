package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	configimpl "grovehub/services/integration/config/impl"
)

func withEnv(name, value string) {
	previous, exists := os.LookupEnv(name)
	Expect(os.Setenv(name, value)).To(Succeed())
	DeferCleanup(func() {
		if exists {
			_ = os.Setenv(name, previous)
		} else {
			_ = os.Unsetenv(name)
		}
	})
}

var _ = Describe("OAuth callback validation", func() {
	cases := []struct {
		provider, envName, uri string
		valid                  bool
	}{
		{"gmail", "GROVEHUB_GMAIL_REDIRECT_URI", "https://example.test/integrations/connections/oauth/gmail/callback", true},
		{"gmail", "GROVEHUB_GMAIL_REDIRECT_URI", "https://example.test/integrations/oauth/gmail/callback", false},
		{"outlook", "GROVEHUB_OUTLOOK_REDIRECT_URI", "https://example.test/integrations/connections/oauth/outlook/callback", true},
		{"outlook", "GROVEHUB_OUTLOOK_REDIRECT_URI", "http://example.test/integrations/connections/oauth/outlook/callback", false},
	}
	for _, testCase := range cases {
		tc := testCase
		It("validates "+tc.provider+" redirect "+tc.uri, func() {
			withEnv(tc.envName, tc.uri)
			var err error
			if tc.provider == "gmail" {
				withEnv("GROVEHUB_GMAIL_CLIENT_ID", "id")
				withEnv("GROVEHUB_GMAIL_CLIENT_SECRET", "secret")
				_, err = configimpl.NewEnvironmentLoader().Gmail()
			} else {
				withEnv("GROVEHUB_OUTLOOK_CLIENT_ID", "id")
				withEnv("GROVEHUB_OUTLOOK_CLIENT_SECRET", "secret")
				_, err = configimpl.NewEnvironmentLoader().Outlook()
			}
			Expect(err == nil).To(Equal(tc.valid))
		})
	}
})

var _=Describe("local Integration authentication bypass",func(){
 It("defaults to disabled",func(){
  withEnv("GROVEHUB_INTEGRATION_DEV_AUTH_BYPASS","false")
  cfg,err:=configimpl.NewEnvironmentLoader().DevAuth()
  Expect(err).NotTo(HaveOccurred())
  Expect(cfg.Enabled).To(BeFalse())
 })
 It("requires valid fixed UUIDv7 identities and local deployment",func(){
  withEnv("GROVEHUB_INTEGRATION_DEV_AUTH_BYPASS","true")
  withEnv("OTEL_DEPLOYMENT_ENVIRONMENT","local")
  withEnv("GROVEHUB_INTEGRATION_DEV_TENANT_ID","018f5d04-315c-75a9-b890-df0d02e787c1")
  withEnv("GROVEHUB_INTEGRATION_DEV_ACTOR_GID","01a1b210-402c-7622-a33c-f82ceba26c7d")
  cfg,err:=configimpl.NewEnvironmentLoader().DevAuth()
  Expect(err).NotTo(HaveOccurred())
  Expect(cfg.Enabled).To(BeTrue())
  Expect(cfg.ActorGID).To(Equal("01a1b210-402c-7622-a33c-f82ceba26c7d"))
 })
 It("rejects production deployment",func(){
  withEnv("GROVEHUB_INTEGRATION_DEV_AUTH_BYPASS","true")
  withEnv("OTEL_DEPLOYMENT_ENVIRONMENT","production")
  withEnv("GROVEHUB_INTEGRATION_DEV_TENANT_ID","018f5d04-315c-75a9-b890-df0d02e787c1")
  withEnv("GROVEHUB_INTEGRATION_DEV_ACTOR_GID","01a1b210-402c-7622-a33c-f82ceba26c7d")
  _,err:=configimpl.NewEnvironmentLoader().DevAuth()
  Expect(err).To(HaveOccurred())
 })
 It("rejects malformed actor identity",func(){
  withEnv("GROVEHUB_INTEGRATION_DEV_AUTH_BYPASS","true")
  withEnv("OTEL_DEPLOYMENT_ENVIRONMENT","local")
  withEnv("GROVEHUB_INTEGRATION_DEV_TENANT_ID","018f5d04-315c-75a9-b890-df0d02e787c1")
  withEnv("GROVEHUB_INTEGRATION_DEV_ACTOR_GID","not-a-uuid")
  _,err:=configimpl.NewEnvironmentLoader().DevAuth()
  Expect(err).To(HaveOccurred())
 })
})
