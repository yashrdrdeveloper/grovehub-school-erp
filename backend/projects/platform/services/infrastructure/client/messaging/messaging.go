// Package messaging exposes Infrastructure's typed public broker client,
// configuration, durable event store and delivery DTOs in one file.
package messaging

import (
    "context"

    inframessaging "grovehub/libs/infrastructure/messaging"
)

type Config = inframessaging.Config
const ProviderNATS = inframessaging.ProviderNATS
func DefaultMessaging() Config { return inframessaging.FromEnv() }

type Event = inframessaging.Event
type EventStore interface {
    Enqueue(ctx context.Context, event inframessaging.Event) error
}

type Message = inframessaging.Message
type MessageHandler = inframessaging.Handler
type PublishOptions = inframessaging.PublishOptions

// @Client messagingClient
type MessagingClient interface {
    // @Proxy service.Publish
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
    // @Proxy service.PublishWithOptions
    PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options inframessaging.PublishOptions) error
    // @Proxy service.Consume
    Consume(ctx context.Context, topic string, handler inframessaging.Handler) error
    // @Proxy service.Health
    Health(ctx context.Context) error
}
