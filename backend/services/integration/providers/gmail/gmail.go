package gmail

import "grovehub/services/integration/providers"

// Gmail provider capabilities are implemented through the connected OAuth client.
type Provider struct{}
func New() providers.Provider { return Provider{} }
func (Provider) Descriptor() providers.Descriptor {
 return providers.Descriptor{
  ID: "GMAIL", DisplayName: "Gmail", AuthType: "OAUTH2",
  Capabilities: []providers.Capability{
   providers.CapabilityMailSend, providers.CapabilityMailRead,
   providers.CapabilityMailHistory, providers.CapabilityMailWatch,
  },
 }
}
