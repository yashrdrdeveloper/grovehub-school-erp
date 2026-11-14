package adapters

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/ThreeDotsLabs/watermill"
	watermillnats "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	infrastructureconfig "grovehub/services/infrastructure/internal/app/config"
	messagingdto "grovehub/services/infrastructure/internal/app/dto/messaging"
    natscontract "grovehub/services/infrastructure/internal/app/externalService/nats"
)

// The vendor adapter is verified against the annotated external service contract.
// Business services depend on MessagingClient, never on NATS directly.

const traceIDMetadataKey = "X-Trace-ID"

var _ natscontract.Service = (*Adapter)(nil)

type Adapter struct {
	url        string
	publisher  message.Publisher
	subscriber message.Subscriber
}

func New(config infrastructureconfig.Messaging) (*Adapter, error) {
	jetStream := watermillnats.JetStreamConfig{
		AutoProvision: true,
		TrackMsgID:    true,
		AckAsync:      false,
		DurablePrefix: config.DurablePrefix,
	}

	publisher, err := watermillnats.NewPublisher(
		watermillnats.PublisherConfig{
			URL:       config.NATSURL,
			JetStream: jetStream,
		},
		watermill.NopLogger{},
	)
	if err != nil {
		return nil, err
	}

	subscriber, err := watermillnats.NewSubscriber(
		watermillnats.SubscriberConfig{
			URL:              config.NATSURL,
			QueueGroupPrefix: config.QueueGroupPrefix,
			CloseTimeout:     config.CloseTimeout,
			AckWaitTimeout:   config.AckWaitTimeout,
			SubscribeTimeout: config.SubscribeTimeout,
			JetStream:        jetStream,
		},
		watermill.NopLogger{},
	)
	if err != nil {
		_ = publisher.Close()
		return nil, err
	}

	return &Adapter{
		url:        config.NATSURL,
		publisher:  publisher,
		subscriber: subscriber,
	}, nil
}

func (adapter *Adapter) Publish(
	ctx context.Context,
	topic string,
	payload []byte,
	metadata map[string]string,
) error {
	msg := message.NewMessage(watermill.NewUUID(), payload)
	msg.SetContext(ctx)

	for key, value := range metadata {
		msg.Metadata.Set(key, value)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Metadata))
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		msg.Metadata.Set(traceIDMetadataKey, spanContext.TraceID().String())
	}

	return adapter.publisher.Publish(topic, msg)
}

func (adapter *Adapter) Consume(
	ctx context.Context,
	topic string,
	handler messagingdto.MessageHandler,
) error {
	messages, err := adapter.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	for msg := range messages {
		messageContext := otel.GetTextMapPropagator().Extract(
			ctx,
			propagation.MapCarrier(msg.Metadata),
		)

		err := handler(messageContext, messagingdto.Message{
			ID:       msg.UUID,
			Payload:  append([]byte(nil), msg.Payload...),
			Metadata: cloneMetadata(msg.Metadata),
		})
		if err != nil {
			msg.Nack()
			continue
		}

		msg.Ack()
	}

	return ctx.Err()
}

func (adapter *Adapter) Health(ctx context.Context) error {
	address, err := natsAddress(adapter.url)
	if err != nil {
		return err
	}

	connection, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}

	return connection.Close()
}

func (adapter *Adapter) Close() error {
	return errors.Join(adapter.subscriber.Close(), adapter.publisher.Close())
}

func cloneMetadata(metadata message.Metadata) map[string]string {
	result := make(map[string]string, len(metadata))
	for key, value := range metadata {
		result[key] = value
	}
	return result
}

func natsAddress(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if parsed.Host == "" {
		return "", errors.New("NATS URL must include a host")
	}

	if strings.Contains(parsed.Host, ":") {
		return parsed.Host, nil
	}

	return net.JoinHostPort(parsed.Host, "4222"), nil
}
