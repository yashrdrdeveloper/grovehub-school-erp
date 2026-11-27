package outlook

import (
    providercontracts "grovehub/services/integration/internal/app/providers/contracts"
	integrationconfig "grovehub/services/integration/internal/app/config"
	"grovehub/services/integration/internal/app/providers"
	impl "grovehub/services/integration/internal/app/externalService/mail/outlook/impl"
)

// Public Outlook provider API. Concrete Graph and OAuth behavior is in impl.
type Client = impl.Client

type Provider = impl.Provider

type Message = impl.Message

type MessagePage = impl.MessagePage

type DeltaPage = impl.DeltaPage

type SendResult = impl.SendResult

type StatusError = impl.StatusError

func New() providers.Provider                         { return impl.New() }

func NewClient(cfg integrationconfig.Outlook) *Client { return impl.NewClient(cfg) }

func RegisterActions(registry providercontracts.ActionRegistry, client *Client) error {
	return impl.RegisterActions(registry, client)
}
