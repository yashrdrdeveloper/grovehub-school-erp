package client

import (
    "context"

    infrastructureservices "grovehub/services/infrastructure/services"
)

type Message = infrastructureservices.Message
type MessageHandler = infrastructureservices.MessageHandler

// @Client messagingClient
type MessagingClient interface {
    // @Proxy service.Publish
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
    // @Proxy service.Consume
    Consume(ctx context.Context, topic string, handler infrastructureservices.MessageHandler) error
    // @Proxy service.Health
    Health(ctx context.Context) error
}
