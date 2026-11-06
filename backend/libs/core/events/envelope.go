// Package events provides generic, typed envelopes for all domain event payloads.
package events

import (
    "context"
    "time"

    coretelemetry "grovehub/libs/core/telemetry"
)

// Envelope defines the public, transport-independent event wire contract.
// Version is the event schema version, NOT the aggregate's optimistic-lock version.
// CorrelationID and CausationID are omitted when their values are not known.
type Envelope[T any] struct {
    EventID       string    `json:"eventId"`
    EventType     string    `json:"eventType"`
    Version       int       `json:"version"`
    AggregateID   string    `json:"aggregateId"`
    OccurredAt    time.Time `json:"occurredAt"`
    CorrelationID string    `json:"correlationId,omitempty"`
    CausationID   string    `json:"causationId,omitempty"`
    Data          T         `json:"data"`
}

// Metadata contains domain-agnostic event identity and schema metadata.
type Metadata struct {
    EventID       string
    EventType     string
    Version       int
    AggregateID   string
    OccurredAt    time.Time
    CorrelationID string
    CausationID   string
}

type correlationKey struct{}
type causationKey struct{}

// WithCorrelationID explicitly propagates the workflow/request correlation ID.
// When no explicit value exists, a valid OpenTelemetry trace ID is used.
func WithCorrelationID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, correlationKey{}, id)
}

// WithCausationID records the originating command ID or triggering event ID.
// For registration, the validated Idempotency-Key identifies that command.
func WithCausationID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, causationKey{}, id)
}

// NewEnvelope builds one versioned event using any domain payload T.
// Metadata overrides context-derived IDs when explicitly specified.
func NewEnvelope[T any](ctx context.Context, metadata Metadata, data T) Envelope[T] {
    correlationID := metadata.CorrelationID
    if correlationID == "" {
        if value, ok := ctx.Value(correlationKey{}).(string); ok {
            correlationID = value
        }
    }
    if correlationID == "" {
        correlationID = coretelemetry.TraceID(ctx)
    }
    causationID := metadata.CausationID
    if causationID == "" {
        if value, ok := ctx.Value(causationKey{}).(string); ok {
            causationID = value
        }
    }
    return Envelope[T]{
        EventID:       metadata.EventID,
        EventType:     metadata.EventType,
        Version:       metadata.Version,
        AggregateID:   metadata.AggregateID,
        OccurredAt:    metadata.OccurredAt,
        CorrelationID: correlationID,
        CausationID:   causationID,
        Data:          data,
    }
}
