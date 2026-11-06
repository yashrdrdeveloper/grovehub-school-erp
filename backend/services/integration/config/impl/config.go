package impl

import (
	"encoding/base64"
	"errors"
	"net/url"
	"os"
	"strings"

	"grovehub/services/integration/config"
)

type environmentLoader struct{}

func NewEnvironmentLoader() config.Loader { return environmentLoader{} }

func (environmentLoader) Gmail() (config.Gmail, error) {
	value := config.Gmail{
		ClientID:     strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_CLIENT_ID")),
		ClientSecret: strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_CLIENT_SECRET")),
		RedirectURI:  strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_REDIRECT_URI")),
	}
	if value.ClientID == "" || value.ClientSecret == "" || value.RedirectURI == "" {
		return config.Gmail{}, errors.New("Gmail OAuth configuration missing")
	}
	if !validCallback(value.RedirectURI, "gmail") {
		return config.Gmail{}, errors.New("Gmail redirect URI must use the connections OAuth callback path over HTTPS")
	}
	return value, nil
}
func (environmentLoader) Outlook() (config.Outlook, error) {
	value := config.Outlook{
		ClientID:     strings.TrimSpace(os.Getenv("GROVEHUB_OUTLOOK_CLIENT_ID")),
		ClientSecret: strings.TrimSpace(os.Getenv("GROVEHUB_OUTLOOK_CLIENT_SECRET")),
		RedirectURI:  strings.TrimSpace(os.Getenv("GROVEHUB_OUTLOOK_REDIRECT_URI")),
		Tenant:       strings.TrimSpace(os.Getenv("GROVEHUB_OUTLOOK_TENANT")),
	}
	if value.Tenant == "" {
		value.Tenant = "common"
	}
	if value.ClientID == "" || value.ClientSecret == "" || value.RedirectURI == "" {
		return config.Outlook{}, errors.New("Outlook OAuth configuration missing")
	}
	if !validCallback(value.RedirectURI, "outlook") {
		return config.Outlook{}, errors.New("Outlook redirect URI must use the connections OAuth callback path over HTTPS")
	}
	if strings.ContainsAny(value.Tenant, "/?#\\ \t\r\n") {
		return config.Outlook{}, errors.New("invalid Outlook tenant")
	}
	return value, nil
}
func (environmentLoader) GmailPush() config.GmailPush {
	return config.GmailPush{
		Topic:          strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_PUBSUB_TOPIC")),
		Audience:       strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_PUBSUB_AUDIENCE")),
		ServiceAccount: strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_PUBSUB_SERVICE_ACCOUNT")),
	}
}
func (environmentLoader) CredentialKey() ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(os.Getenv("GROVEHUB_INTEGRATION_ENCRYPTION_KEY"))
	if err != nil || len(key) != 32 {
		return nil, errors.New("integration encryption key must be base64-encoded 32 bytes")
	}
	return key, nil
}
func secureRedirect(uri string) bool {
	return strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://localhost:")
}

// validCallback prevents an obsolete provider redirect from being configured
// while the HTTP router only exposes the new canonical connections endpoint.
func validCallback(raw, provider string) bool {
	if !secureRedirect(raw) {
		return false
	}
	uri, err := url.Parse(raw)
	if err != nil || uri.User != nil || uri.Hostname() == "" || uri.RawQuery != "" || uri.Fragment != "" {
		return false
	}
	return uri.Path == "/integrations/connections/oauth/"+provider+"/callback"
}
