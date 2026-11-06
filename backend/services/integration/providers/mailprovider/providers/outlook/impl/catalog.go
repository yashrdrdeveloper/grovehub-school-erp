package impl

import "grovehub/services/integration/providers"

// Outlook provider descriptor; supported action strategies are registered separately.
type Provider struct{}
func New() providers.Provider { return Provider{} }
func (Provider) Descriptor() providers.Descriptor {
 return providers.Descriptor{
  ID: "OUTLOOK", DisplayName: "Outlook", AuthType: "OAUTH2",
  Capabilities: []providers.Capability{providers.CapabilityMailSend,providers.CapabilityMailRead,providers.CapabilityMailDelta,providers.CapabilityMailAttachment}, // mail.history and watch still unsupported
 }
}
