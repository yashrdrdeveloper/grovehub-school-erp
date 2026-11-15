package publish

import (
    "context"
    coremessaging "grovehub/libs/core/events/messaging"
    "grovehub/services/infrastructure/internal/app/services/messaging/operations/publish/strategies"
)

type PublishOperation struct { strategy strategies.Strategy }
func New(strategy strategies.Strategy) *PublishOperation { return &PublishOperation{strategy:strategy} }
func (operation *PublishOperation) Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error {
    return operation.strategy.Publish(ctx, topic, payload, metadata)
}
func (operation *PublishOperation) PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options coremessaging.PublishOptions) error {
    return operation.strategy.PublishWithOptions(ctx, topic, payload, metadata, options)
}
