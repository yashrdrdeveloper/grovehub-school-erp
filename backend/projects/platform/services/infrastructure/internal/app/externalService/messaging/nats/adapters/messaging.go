package adapters

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"
	"regexp"
	inframessaging "grovehub/services/infrastructure/client/messaging"

	"github.com/ThreeDotsLabs/watermill"
	watermillnats "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	infrastructureconfig "grovehub/services/infrastructure/internal/app/config/messaging"
	messagingdto "grovehub/services/infrastructure/internal/app/dto/messaging"
    natscontract "grovehub/services/infrastructure/internal/app/externalService/messaging/nats"
    messagingstrategies "grovehub/services/infrastructure/internal/app/services/messaging/strategies"
)

// The vendor adapter is verified against the annotated external service contract.
// Business services depend on MessagingClient, never on NATS directly.

const traceIDMetadataKey = "X-Trace-ID"

var consumerGroupPattern = regexp.MustCompile("^[a-z][a-z0-9-]{0,31}$")

var _ natscontract.Service = (*Adapter)(nil)
var _ messagingstrategies.Strategy = (*Adapter)(nil)

type Adapter struct {
	url        string
	publisher  message.Publisher
	subscriber message.Subscriber
	subscriberConfig watermillnats.SubscriberConfig
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

	subscriberConfig := watermillnats.SubscriberConfig{
			URL:              config.NATSURL,
			QueueGroupPrefix: config.QueueGroupPrefix,
			CloseTimeout:     config.CloseTimeout,
			AckWaitTimeout:   config.AckWaitTimeout,
			SubscribeTimeout: config.SubscribeTimeout,
			JetStream:        jetStream,
	}
	subscriber, err := watermillnats.NewSubscriber(subscriberConfig,watermill.NopLogger{})
	if err != nil {
		_ = publisher.Close()
		return nil, err
	}

	return &Adapter{
		url:        config.NATSURL,
		publisher:  publisher,
		subscriber: subscriber,
		subscriberConfig:subscriberConfig,
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
	subscriber := adapter.subscriber
	if group := inframessaging.ConsumerGroup(ctx); group != "" {
		if !consumerGroupPattern.MatchString(group) {
			return errors.New("invalid consumer group")
		}
		config := adapter.subscriberConfig
		config.QueueGroupPrefix += "-" + group
		config.JetStream.DurablePrefix += "-" + group
		var err error
		subscriber, err = watermillnats.NewSubscriber(config, watermill.NopLogger{})
		if err != nil {
			return err
		}
		defer subscriber.Close()
	}
	messages, err := subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	inframessaging.NotifySubscription(ctx,true)
	defer inframessaging.NotifySubscription(ctx,false)

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
