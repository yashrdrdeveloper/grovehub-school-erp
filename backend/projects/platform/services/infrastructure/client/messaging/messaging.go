// Package messaging exposes Infrastructure's transport-neutral broker client
// and transactional event-store contracts.
package messaging

import (
	"context"

	config "grovehub/services/infrastructure/internal/app/config/messaging"
	"grovehub/services/infrastructure/internal/app/client/messaging/helpers"
)

type Config = config.Messaging

const ProviderNATS = config.ProviderNATS

func DefaultMessaging() Config { return config.FromEnv() }
func FromEnv() Config { return config.FromEnv() }

// Event is a transport-neutral pending outbox message.
type Event struct {
	EventID string
	Topic string
	Payload []byte
	Metadata map[string]string
}

type EventStore interface {
	Enqueue(ctx context.Context, event Event) error
}

// Message represents a broker delivery independent of provider types.
type Message struct {
	ID string
	Payload []byte
	Metadata map[string]string
}

type Handler func(context.Context, Message) error
type MessageHandler = Handler

// PublishOptions selects transactional outbox persistence for a message.
type PublishOptions struct {
	Persist bool
}

type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
}

// WithConsumerGroup selects an independent logical subscriber.
func WithConsumerGroup(ctx context.Context, group string) context.Context {
	return helpers.WithConsumerGroup(ctx, group)
}

func ConsumerGroup(ctx context.Context) string {
	return helpers.ConsumerGroup(ctx)
}

// WithSubscriptionObserver reports subscription readiness and termination.
// The observer must be safe for concurrent use.
func WithSubscriptionObserver(ctx context.Context, observer func(bool)) context.Context {
	return helpers.WithSubscriptionObserver(ctx, observer)
}

func NotifySubscription(ctx context.Context, ready bool) {
	helpers.NotifySubscription(ctx, ready)
}

// @Client messagingClient
type MessagingClient interface {
	// @Proxy service.Publish
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
	// @Proxy service.PublishWithOptions
	PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options PublishOptions) error
	// @Proxy service.Consume
	Consume(ctx context.Context, topic string, handler Handler) error
	// @Proxy service.Health
	Health(ctx context.Context) error
}
