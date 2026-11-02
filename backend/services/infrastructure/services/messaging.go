package services

import "context"

type Message struct {
	ID       string
	Payload  []byte
	Metadata map[string]string
}

type MessageHandler func(ctx context.Context, message Message) error

type MessagingService interface {
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
	Consume(ctx context.Context, topic string, handler MessageHandler) error
	Health(ctx context.Context) error
	Close() error
}
