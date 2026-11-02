package messaging

import (
	"context"

	infrastructureservices "grovehub/services/infrastructure/services"
	natsadapters "grovehub/services/infrastructure/externalService/nats/adapters"
)

type natsStrategy struct {
	adapter *natsadapters.Adapter
}

func NewNATSStrategy(adapter *natsadapters.Adapter) Strategy {
	return &natsStrategy{adapter: adapter}
}

func (strategy *natsStrategy) Publish(
	ctx context.Context,
	topic string,
	payload []byte,
	metadata map[string]string,
) error {
	return strategy.adapter.Publish(ctx, topic, payload, metadata)
}

func (strategy *natsStrategy) Consume(
	ctx context.Context,
	topic string,
	handler infrastructureservices.MessageHandler,
) error {
	return strategy.adapter.Consume(ctx, topic, handler)
}

func (strategy *natsStrategy) Health(ctx context.Context) error {
	return strategy.adapter.Health(ctx)
}

func (strategy *natsStrategy) Close() error {
	return strategy.adapter.Close()
}
