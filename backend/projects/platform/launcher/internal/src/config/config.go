package config

import "grovehub/libs/core/annotation"

// Config is decoded once at startup; an invalid tenant-required flag fails closed.
// @Configuration
type Config struct {
	// @Env
	// @Default
	Address string `env:"GROOVEHUB_HTTP_ADDR" default:"127.0.0.1:8080"`
	// @Env
	AllowedOrigin string `env:"IDENTITY_ALLOWED_ORIGIN"`
	// @Env
	Workload string `env:"PLATFORM_WORKLOAD_CONFIG"`
	// @Env
	TenantOrigin string `env:"PLATFORM_TENANT_ACCESS_ORIGIN"`
	// @Env
	TenantKeyID string `env:"PLATFORM_TENANT_ACCESS_KEY_ID"`
	// @Env
	// @Secret
	TenantPrivateKey string `env:"PLATFORM_TENANT_ACCESS_PRIVATE_KEY" secret:"true"`
	// @Env
	// @Default
	TenantRequired bool `env:"PLATFORM_TENANT_ACCESS_REQUIRED" default:"false"`
}

func FromEnv() (Config, error) {
	return annotation.ConfigurationFromEnv[Config]()
}
