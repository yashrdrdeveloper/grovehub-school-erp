package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/launcher/internal/src/config"
)

// @Testing suite=LauncherConfig
func ConfigSpecs() {
	Describe("FromEnv", func() {
		It("uses the default listener and reads every environment binding", func() {
			t := GinkgoT()
			t.Setenv("GROOVEHUB_HTTP_ADDR", "")
			t.Setenv("IDENTITY_ALLOWED_ORIGIN", "https://school.example")
			t.Setenv("PLATFORM_WORKLOAD_CONFIG", "{}")
			t.Setenv("PLATFORM_TENANT_ACCESS_ORIGIN", "https://platform.example")
			t.Setenv("PLATFORM_TENANT_ACCESS_KEY_ID", "key")
			t.Setenv("PLATFORM_TENANT_ACCESS_PRIVATE_KEY", "test-key")
			t.Setenv("PLATFORM_TENANT_ACCESS_REQUIRED", "true")
			value, err := config.FromEnv()
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal(config.Config{
				Address: "127.0.0.1:8080", AllowedOrigin: "https://school.example",
				Workload: "{}", TenantOrigin: "https://platform.example",
				TenantKeyID: "key", TenantPrivateKey: "test-key", TenantRequired: true,
			}))
		})
		It("preserves an explicit address and loads a false flag", func() {
			t := GinkgoT()
			t.Setenv("GROOVEHUB_HTTP_ADDR", "127.0.0.1:9090")
			t.Setenv("PLATFORM_TENANT_ACCESS_REQUIRED", "false")
			value, err := config.FromEnv()
			Expect(err).NotTo(HaveOccurred())
			Expect(value.Address).To(Equal("127.0.0.1:9090"))
			Expect(value.TenantRequired).To(BeFalse())
		})
		It("rejects malformed required-tenant flags instead of disabling the policy", func() {
			GinkgoT().Setenv("PLATFORM_TENANT_ACCESS_REQUIRED", "invalid")
			_, err := config.FromEnv()
			Expect(err).To(MatchError(ContainSubstring("PLATFORM_TENANT_ACCESS_REQUIRED")))
		})
	})
}
