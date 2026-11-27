package impl

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    mailconstants "grovehub/services/integration/internal/app/providers/mailprovider/constants"
    "grovehub/services/integration/internal/app/providers"
)
// Outlook provider descriptor; supported action strategies are registered separately.
type Provider struct{}

func New() providers.Provider { return Provider{} }

func (Provider) Descriptor() providerdto.Descriptor {
	return providerdto.Descriptor{
		ID: "OUTLOOK", DisplayName: "Outlook", AuthType: "OAUTH2",
		Capabilities: []providerdto.Capability{mailconstants.CapabilityMailSend, mailconstants.CapabilityMailRead, mailconstants.CapabilityMailDelta, mailconstants.CapabilityMailAttachment}, // mail.history and watch still unsupported
	}
}
