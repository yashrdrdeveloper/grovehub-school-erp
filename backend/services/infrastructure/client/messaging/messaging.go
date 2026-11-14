// Package messaging exposes Infrastructure's transport-neutral messaging contract.
package messaging

import (
    "context"

    messagingdto "grovehub/services/infrastructure/internal/app/dto/messaging"
)

type Message = messagingdto.Message
type MessageHandler = messagingdto.MessageHandler

// @Client messagingClient
type MessagingClient interface {
    // @Proxy service.Publish
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
    // @Proxy service.Consume
    Consume(ctx context.Context, topic string, handler messagingdto.MessageHandler) error
    // @Proxy service.Health
    Health(ctx context.Context) error
}
