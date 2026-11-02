package infrastructure

import (
	infrastructureclient "grovehub/services/infrastructure/client"
	clientimpl "grovehub/services/infrastructure/client/impl"
	infrastructureconfig "grovehub/services/infrastructure/config"
	natsadapters "grovehub/services/infrastructure/externalService/nats/adapters"
	serviceimpl "grovehub/services/infrastructure/services/impl"
	"grovehub/services/infrastructure/services/strategy"
	messagingstrategy "grovehub/services/infrastructure/services/strategy/messaging"
)

type Module interface {
	MessagingClient() infrastructureclient.MessagingClient
	Close() error
}

type module struct {
	messagingClient  infrastructureclient.MessagingClient
	messagingService interface{ Close() error }
}

func New(messagingConfig infrastructureconfig.Messaging) (Module, error) {
	natsAdapter, err := natsadapters.New(messagingConfig)
	if err != nil {
		return nil, err
	}

	registry := strategy.NewRegistry[messagingstrategy.Strategy]()
	registry.Register(infrastructureconfig.ProviderNATS, messagingstrategy.NewNATSStrategy(natsAdapter))

	messagingService := serviceimpl.NewMessagingService(messagingConfig.Provider, registry)

	return &module{
		messagingClient:  clientimpl.NewMessagingClient(messagingService),
		messagingService: messagingService,
	}, nil
}

func (module *module) MessagingClient() infrastructureclient.MessagingClient {
	return module.messagingClient
}

func (module *module) Close() error {
	return module.messagingService.Close()
}
