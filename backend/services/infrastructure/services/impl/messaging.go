package impl

import (
	"context"

	coreerrors "grovehub/libs/core/errors"
	errorcode "grovehub/services/infrastructure/constants/error"
	infrastructureservices "grovehub/services/infrastructure/services"
	"grovehub/services/infrastructure/services/strategy"
	messagingstrategy "grovehub/services/infrastructure/services/strategy/messaging"
)

type messagingService struct {
	provider string
	registry *strategy.Registry[messagingstrategy.Strategy]
}

func NewMessagingService(
	provider string,
	registry *strategy.Registry[messagingstrategy.Strategy],
) infrastructureservices.MessagingService {
	return &messagingService{provider: provider, registry: registry}
}

func (service *messagingService) Publish(
	ctx context.Context,
	topic string,
	payload []byte,
	metadata map[string]string,
) error {
	provider, err := service.registry.Resolve(service.provider)
	if err != nil {
		return coreerrors.Wrap(errorcode.MessagingProviderUnavailable, err)
	}

	if err := provider.Publish(ctx, topic, payload, metadata); err != nil {
		return coreerrors.Wrap(errorcode.MessagingPublishFailed, err)
	}

	return nil
}

func (service *messagingService) Consume(
	ctx context.Context,
	topic string,
	handler infrastructureservices.MessageHandler,
) error {
	provider, err := service.registry.Resolve(service.provider)
	if err != nil {
		return coreerrors.Wrap(errorcode.MessagingProviderUnavailable, err)
	}

	if err := provider.Consume(ctx, topic, handler); err != nil {
		return coreerrors.Wrap(errorcode.MessagingConsumeFailed, err)
	}

	return nil
}

func (service *messagingService) Health(ctx context.Context) error {
	provider, err := service.registry.Resolve(service.provider)
	if err != nil {
		return coreerrors.Wrap(errorcode.MessagingProviderUnavailable, err)
	}

	return provider.Health(ctx)
}

func (service *messagingService) Close() error {
	provider, err := service.registry.Resolve(service.provider)
	if err != nil {
		return coreerrors.Wrap(errorcode.MessagingProviderUnavailable, err)
	}

	return provider.Close()
}
