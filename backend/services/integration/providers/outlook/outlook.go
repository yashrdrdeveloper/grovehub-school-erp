package outlook

import "grovehub/services/integration/providers"

// Descriptor-only Phase 1 provider. Microsoft Graph/OAuth are Phase 3.
type Provider struct{}
func New() providers.Provider { return Provider{} }
func (Provider) Descriptor() providers.Descriptor {
 return providers.Descriptor{
  ID: "OUTLOOK", DisplayName: "Outlook", AuthType: "OAUTH2",
  Capabilities: []providers.Capability{
   providers.CapabilityMailSend, providers.CapabilityMailRead,
   providers.CapabilityMailSync, providers.CapabilityMailWatch,
  },
 }
}
