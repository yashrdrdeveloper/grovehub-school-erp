// Package messaging holds transport-neutral event messaging contracts.
package messaging

import "context"

// Publisher is implemented by Infrastructure MessagingClient.
type Publisher interface {
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
}
