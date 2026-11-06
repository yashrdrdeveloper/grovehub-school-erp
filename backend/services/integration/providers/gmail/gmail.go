package gmail

import "grovehub/services/integration/providers"

// Descriptor-only Phase 1 provider. OAuth/send/sync are NOT implemented.
type Provider struct{}
func New() providers.Provider { return Provider{} }
func (Provider) Descriptor() providers.Descriptor {
 return providers.Descriptor{
  ID: "GMAIL", DisplayName: "Gmail", AuthType: "OAUTH2",
  Capabilities: []providers.Capability{
   providers.CapabilityMailSend, providers.CapabilityMailRead,
   providers.CapabilityMailSync, providers.CapabilityMailWatch,
  },
 }
}
