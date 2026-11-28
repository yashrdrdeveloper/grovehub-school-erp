// Package events provides generic, transport-independent domain event envelopes.
package events

import (
	"context"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/causation"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/correlation"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/envelope"
)

// Envelope is the public event wire contract.
type Envelope[T any] = contracts.Envelope[T]

// Metadata carries event identity and schema metadata.
type Metadata = contracts.Metadata

// WithCorrelationID explicitly propagates the workflow/request correlation ID.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return correlation.WithID(ctx, id)
}

// WithCausationID records the originating command or triggering event ID.
func WithCausationID(ctx context.Context, id string) context.Context {
	return causation.WithID(ctx, id)
}

// NewEnvelope builds a versioned event. Explicit metadata overrides context IDs.
func NewEnvelope[T any](ctx context.Context, metadata Metadata, data T) Envelope[T] {
	return envelope.New(ctx, metadata, data)
}
