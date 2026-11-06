package outlook

import "grovehub/services/integration/providers"

// Descriptor-only provider; do not advertise unsupported operations.
type Provider struct{}
func New() providers.Provider { return Provider{} }
func (Provider) Descriptor() providers.Descriptor {
 return providers.Descriptor{
  ID: "OUTLOOK", DisplayName: "Outlook", AuthType: "OAUTH2",
  Capabilities: []providers.Capability{providers.CapabilityMailSend,providers.CapabilityMailRead,providers.CapabilityMailDelta,providers.CapabilityMailAttachment}, // mail.history and watch still unsupported
 }
}
