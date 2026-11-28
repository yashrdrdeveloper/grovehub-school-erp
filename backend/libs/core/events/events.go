// Package events exposes transport-independent event envelopes and execution rules.
package events

import (
	"context"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/causation"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/correlation"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/envelope"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/pipeline"
)

// Envelope is the public event wire contract.
type Envelope[T any] = contracts.Envelope[T]

// Metadata carries event identity and schema metadata.
type Metadata = contracts.Metadata

// Rules composes eligibility, construction and delivery using a consumer-owned message.
type Rules[T any, M any] = contracts.Rules[T, M]

// ErrIncompleteRules indicates that a required pipeline callback is missing.
var ErrIncompleteRules = contracts.ErrIncompleteRules

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

// Run executes event eligibility, construction and enqueue in order.
func Run[T any, M any](ctx context.Context, value T, rules Rules[T, M]) error {
	return pipeline.Run(ctx, value, rules)
}
