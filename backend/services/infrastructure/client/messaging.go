package client

import (
	"context"

	infrastructureservices "grovehub/services/infrastructure/services"
)

type Message = infrastructureservices.Message
type MessageHandler = infrastructureservices.MessageHandler

type MessagingClient interface {
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
	Consume(ctx context.Context, topic string, handler MessageHandler) error
	Health(ctx context.Context) error
}
