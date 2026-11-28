package strategies

import (
    "context"
    inframessaging "grovehub/services/infrastructure/client/messaging"
)

// @Strategy
type Strategy interface {
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
    PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options inframessaging.PublishOptions) error
}
