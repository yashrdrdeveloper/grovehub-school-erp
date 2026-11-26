// Package producer is the public Audit event publication contract.
package producer

import "context"

// Event contains a versioned domain envelope. Audit validates the schema and
// owns subject selection and transactional outbox publication.
type Event struct {
	Type string
	Payload []byte
}

// @Client producerClient
type Client interface {
	// @Proxy service.Produce
	Produce(ctx context.Context, event Event) error
}
