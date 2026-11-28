package messaging

import "context"

// PublishOptions selects transactional outbox persistence on a per-message basis.
// Persist:false sends directly to the provider; Persist:true requires the
// business UnitOfWork transaction context for atomic state/event commits.
type PublishOptions struct {
    Persist bool
}

type Publisher interface {
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
}
