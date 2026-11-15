package strategies

import (
    "context"
    coremessaging "grovehub/libs/core/events/messaging"
)

// @Strategy
type Strategy interface {
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
    PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options coremessaging.PublishOptions) error
}
