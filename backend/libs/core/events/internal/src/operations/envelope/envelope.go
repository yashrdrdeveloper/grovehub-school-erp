package envelope

import (
	"context"

	"grovehub/libs/core/events/internal/src/contracts"
	"grovehub/libs/core/events/internal/src/operations/causation"
	"grovehub/libs/core/events/internal/src/operations/correlation"
)

func New[T any](ctx context.Context, metadata contracts.Metadata, data T) contracts.Envelope[T] {
	return contracts.Envelope[T]{
		EventID:       metadata.EventID,
		EventType:     metadata.EventType,
		Version:       metadata.Version,
		AggregateID:   metadata.AggregateID,
		OccurredAt:    metadata.OccurredAt,
		CorrelationID: correlation.Resolve(ctx, metadata.CorrelationID),
		CausationID:   causation.Resolve(ctx, metadata.CausationID),
		Data:          data,
	}
}
