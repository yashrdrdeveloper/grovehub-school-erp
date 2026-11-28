package impl_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"grovehub/services/infrastructure/client/messaging"
	clientimpl "grovehub/services/infrastructure/internal/app/client/messaging/impl"
)

type service struct {
	ctx context.Context
	topic string
	payload []byte
	metadata map[string]string
	options messaging.PublishOptions
	err error
}

func (s *service) Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error {
	s.ctx, s.topic, s.payload, s.metadata = ctx, topic, payload, metadata
	return s.err
}

func (s *service) PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options messaging.PublishOptions) error {
	s.options = options
	return s.Publish(ctx, topic, payload, metadata)
}

func (s *service) Consume(ctx context.Context, topic string, handler messaging.Handler) error {
	s.ctx, s.topic = ctx, topic
	if err := handler(ctx, messaging.Message{ID: "message"}); err != nil {
		return err
	}
	return s.err
}

func (s *service) Health(ctx context.Context) error {
	s.ctx = ctx
	return s.err
}

func (s *service) Close() error {
	return nil
}

func TestMessagingClient(t *testing.T) {
	ctx := messaging.WithConsumerGroup(context.Background(), "group")
	payload := []byte("payload")
	metadata := map[string]string{"correlationId": "correlation"}
	for _, failure := range []error{nil, errors.New("broker failed")} {
		s := &service{err: failure}
		client := clientimpl.NewMessagingClient(s)
		if err := client.Publish(ctx, "topic", payload, metadata); err != failure || s.ctx != ctx || s.topic != "topic" || !reflect.DeepEqual(s.payload, payload) || !reflect.DeepEqual(s.metadata, metadata) {
			t.Fatal("publication forwarding failed")
		}
		if err := client.PublishWithOptions(ctx, "topic", payload, metadata, messaging.PublishOptions{Persist: true}); err != failure || !s.options.Persist {
			t.Fatal("publication options not forwarded")
		}
		if err := client.Consume(ctx, "topic", func(got context.Context, message messaging.Message) error {
			if got != ctx || message.ID != "message" {
				t.Fatal("delivery context or message changed")
			}
			return nil
		}); err != failure {
			t.Fatal("consume error not forwarded")
		}
		if err := client.Health(ctx); err != failure || s.ctx != ctx {
			t.Fatal("health forwarding failed")
		}
	}
}
