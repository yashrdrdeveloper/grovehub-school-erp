package outlook

import (
	providercontracts "grovehub/services/integration/internal/app/providers/contracts"
	providerdto "grovehub/services/integration/internal/app/providers/dto"
	"grovehub/services/integration/internal/app/providers/mailprovider"
	mailcontracts "grovehub/services/integration/internal/app/providers/mailprovider/contracts"
)

// @Provider id=OUTLOOK category=mail
// @Strategy
type strategy struct {
	client *Client
}

var _ mailprovider.Strategy = (*strategy)(nil)

func NewStrategy(client *Client) mailprovider.Strategy {
	return &strategy{client: client}
}

func (provider *strategy) Descriptor() providerdto.Descriptor {
	return New().Descriptor()
}

func (provider *strategy) Register(actions providercontracts.ActionRegistry, oauth mailcontracts.OAuthRegistry) error {
	if provider.client == nil {
		return nil
	}
	if err := RegisterActions(actions, provider.client); err != nil {
		return err
	}
	return oauth.Register("OUTLOOK", provider.client)
}
