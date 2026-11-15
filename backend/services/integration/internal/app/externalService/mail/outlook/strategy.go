package outlook

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    providercontracts "grovehub/services/integration/internal/app/providers/contracts"
    mailcontracts "grovehub/services/integration/internal/app/providers/mailprovider/contracts"
	"grovehub/services/integration/internal/app/providers/mailprovider"
)

// @Provider id=OUTLOOK category=mail
type strategy struct{ client *Client }

func NewStrategy(client *Client) mailprovider.Strategy { return &strategy{client: client} }

func (s *strategy) Descriptor() providerdto.Descriptor   { return New().Descriptor() }

func (s *strategy) Register(actions providercontracts.ActionRegistry, oauth mailcontracts.OAuthRegistry) error {
	if s.client == nil {
		return nil
	}
	if err := RegisterActions(actions, s.client); err != nil {
		return err
	}
	return oauth.Register("OUTLOOK", s.client)
}
