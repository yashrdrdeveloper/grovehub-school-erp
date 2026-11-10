package config

// Config contracts are provider-neutral. Environment parsing belongs in
// config/impl and uses the typed @Configuration startup loader.

// @Configuration
type Gmail struct {
    ClientID     string `env:"GROVEHUB_GMAIL_CLIENT_ID" required:"true"`
    ClientSecret string `env:"GROVEHUB_GMAIL_CLIENT_SECRET" required:"true" secret:"true"`
    RedirectURI  string `env:"GROVEHUB_GMAIL_REDIRECT_URI" required:"true"`
}

// @Configuration
type Outlook struct {
    ClientID     string `env:"GROVEHUB_OUTLOOK_CLIENT_ID" required:"true"`
    ClientSecret string `env:"GROVEHUB_OUTLOOK_CLIENT_SECRET" required:"true" secret:"true"`
    RedirectURI  string `env:"GROVEHUB_OUTLOOK_REDIRECT_URI" required:"true"`
    Tenant       string `env:"GROVEHUB_OUTLOOK_TENANT" default:"common"`
}

// @Configuration
type GmailPush struct {
    Topic          string `env:"GROVEHUB_GMAIL_PUBSUB_TOPIC"`
    Audience       string `env:"GROVEHUB_GMAIL_PUBSUB_AUDIENCE"`
    ServiceAccount string `env:"GROVEHUB_GMAIL_PUBSUB_SERVICE_ACCOUNT"`
}
func (c GmailPush) Enabled() bool {
    return c.Topic != "" && c.Audience != "" && c.ServiceAccount != ""
}

// DevAuth is an explicitly opt-in principal for local Integration testing.
// Never derive the principal from caller-controlled request headers.
type DevAuth struct {
    Enabled bool
    TenantID string
    ActorGID string
}

type Loader interface {
    Gmail() (Gmail,error)
    Outlook() (Outlook,error)
    GmailPush() GmailPush
    CredentialKey() ([]byte,error)
}
