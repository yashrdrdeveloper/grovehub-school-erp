package impl

import (
	"context"

	infrastructureclient "grovehub/services/infrastructure/client"
	infrastructureservices "grovehub/services/infrastructure/services"
)

type messagingClient struct {
	service infrastructureservices.MessagingService
}

func NewMessagingClient(service infrastructureservices.MessagingService) infrastructureclient.MessagingClient {
	return &messagingClient{service: service}
}

func (client *messagingClient) Publish(
	ctx context.Context,
	topic string,
	payload []byte,
	metadata map[string]string,
) error {
	return client.service.Publish(ctx, topic, payload, metadata)
}

func (client *messagingClient) Consume(
	ctx context.Context,
	topic string,
	handler infrastructureclient.MessageHandler,
) error {
	return client.service.Consume(ctx, topic, handler)
}

func (client *messagingClient) Health(ctx context.Context) error {
	return client.service.Health(ctx)
}
