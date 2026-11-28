package impl

import (
    "context"
    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"

    client "grovehub/services/infrastructure/client/messaging"
    errorcode "grovehub/services/infrastructure/internal/app/constants/error"
    "grovehub/services/infrastructure/internal/app/services/messaging/operations/publish/helpers"
    publishstrategies "grovehub/services/infrastructure/internal/app/services/messaging/operations/publish/strategies"
    providerstrategies "grovehub/services/infrastructure/internal/app/services/messaging/strategies"
)

type publishStrategy struct {
    provider providerstrategies.Strategy
    store client.EventStore
}
var _ publishstrategies.Strategy = (*publishStrategy)(nil)

func New(provider providerstrategies.Strategy, store client.EventStore) publishstrategies.Strategy {
    return &publishStrategy{provider:provider, store:store}
}
func (strategy *publishStrategy) Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error {
    if err := strategy.provider.Publish(ctx, topic, payload, metadata); err != nil {
        return coreerrors.Wrap(errorcode.MessagingPublishFailed, err)
    }
    return nil
}
func (strategy *publishStrategy) PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options client.PublishOptions) error {
    if !options.Persist { return strategy.Publish(ctx, topic, payload, metadata) }
    if strategy.store == nil { return errorcode.ErrInvalidStore }
    // The EventStore uses the business UnitOfWork transaction bound to ctx.
    return strategy.store.Enqueue(ctx, helpers.NewMessage(topic, payload, metadata))
}
