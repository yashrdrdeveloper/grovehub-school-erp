package providers

import "context"

// Capability is a stable logical operation; providers declare only what they support.
type Capability string

const (
 CapabilityMailSend Capability = "mail.send"
 CapabilityMailRead Capability = "mail.read"
 CapabilityMailAttachment Capability = "mail.attachment.get"
 CapabilityMailHistory Capability = "mail.history"
 CapabilityMailDelta Capability = "mail.delta"
 CapabilityMailWatch Capability = "mail.watch"
)

// Descriptor describes a provider type, not a live authenticated connection.
type Descriptor struct {
 ID string `json:"id"`
 DisplayName string `json:"displayName"`
 Capabilities []Capability `json:"capabilities"`
 AuthType string `json:"authType"`
}

// Provider is a registration/catalog contract; OAuth and runtime actions are
// separate optional capability interfaces added when adapters exist.
type Provider interface { Descriptor() Descriptor }

// OAuthConnector is optional; registering a provider never implies OAuth works.
type OAuthConnector interface {
 StartOAuth(ctx context.Context, connectionID, redirectURI string) (authorizationURL string, error)
 CompleteOAuth(ctx context.Context, connectionID, code, state string) (credentialRef string, error)
}
