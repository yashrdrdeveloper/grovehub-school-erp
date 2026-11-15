// Package messaging holds transport-neutral event messaging contracts.
package messaging

import "context"

// PublishOptions controls whether a message is recorded in PostgreSQL before
// asynchronous delivery. The default (Persist: false) publishes to the broker
// directly. Persistence must occur in the business UnitOfWork to be atomic with
// a business mutation.
type PublishOptions struct {
	Persist bool
}

// Publisher is implemented by Infrastructure MessagingClient.
type Publisher interface {
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
}
