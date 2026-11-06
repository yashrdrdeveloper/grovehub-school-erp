package config

// Config contracts are provider-neutral from the composition root's
// perspective. Runtime environment parsing belongs in config/impl.
type Gmail struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}
type Outlook struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Tenant       string
}
type GmailPush struct {
	Topic          string
	Audience       string
	ServiceAccount string
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
	Gmail() (Gmail, error)
	Outlook() (Outlook, error)
	GmailPush() GmailPush
	CredentialKey() ([]byte, error)
 DevAuth() (DevAuth,error)
}
