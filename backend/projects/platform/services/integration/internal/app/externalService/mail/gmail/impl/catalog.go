package impl

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    mailconstants "grovehub/services/integration/internal/app/providers/mailprovider/constants"
    "grovehub/services/integration/internal/app/providers"
)
// Gmail provider capabilities are implemented through the connected OAuth client.
type Provider struct{}

func New() providers.Provider { return Provider{} }

func (Provider) Descriptor() providerdto.Descriptor {
	return providerdto.Descriptor{
		ID: "GMAIL", DisplayName: "Gmail", AuthType: "OAUTH2",
		Capabilities: []providerdto.Capability{
			mailconstants.CapabilityMailSend, mailconstants.CapabilityMailRead,
			mailconstants.CapabilityMailHistory, mailconstants.CapabilityMailWatch, mailconstants.CapabilityMailAttachment,
		},
	}
}
