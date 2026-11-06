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
