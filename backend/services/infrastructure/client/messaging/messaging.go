// Package messaging exposes Infrastructure's transport-neutral messaging contract.
package messaging

import (
    "context"

    coremessaging "grovehub/libs/core/events/messaging"
)

type Message = coremessaging.Message
type MessageHandler = coremessaging.Handler

// @Client messagingClient
type MessagingClient interface {
    // @Proxy service.Publish
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
    // @Proxy service.Consume
    Consume(ctx context.Context, topic string, handler coremessaging.Handler) error
    // @Proxy service.Health
    Health(ctx context.Context) error
}
