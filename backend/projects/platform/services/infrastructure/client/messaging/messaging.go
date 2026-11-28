// Package messaging exposes Infrastructure's transport-neutral broker client
// and transactional event-store contracts.
package messaging

import (
	"context"

	"grovehub/services/infrastructure/client/messaging/constants"
	"grovehub/services/infrastructure/internal/app/client/messaging/helpers"
	config "grovehub/services/infrastructure/internal/app/config/messaging"
	dto "grovehub/services/infrastructure/internal/app/dto/messaging"
)

type (
	Config = config.Messaging
	Event = dto.Event
	Message = dto.Message
	Handler = dto.Handler
	MessageHandler = dto.MessageHandler
	PublishOptions = dto.PublishOptions
)

const ProviderNATS = config.ProviderNATS

var ErrReplayUnavailable = constants.ErrReplayUnavailable

// @Client messagingClient
type Client interface {
	// @Proxy service.Publish
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error

	// @Proxy service.PublishWithOptions
	PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options dto.PublishOptions) error

	// @Proxy service.Consume
	Consume(ctx context.Context, topic string, handler dto.Handler) error

	// @Proxy service.Health
	Health(ctx context.Context) error
}

// MessagingClient preserves the existing public type name.
type MessagingClient = Client

type EventStore interface {
	Enqueue(ctx context.Context, event Event) error
}

type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
}

// ReplaySource is an operator-only, read-only view of retained published events.
type ReplaySource interface {
	Published(ctx context.Context, eventID string) (Event, error)
}

func DefaultMessaging() Config {
	return config.FromEnv()
}

func FromEnv() Config {
	return config.FromEnv()
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
