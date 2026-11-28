package strategies

import (
    "context"
    inframessaging "grovehub/libs/infrastructure/messaging"
)

// @Strategy
type Strategy interface {
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
    PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options inframessaging.PublishOptions) error
}
