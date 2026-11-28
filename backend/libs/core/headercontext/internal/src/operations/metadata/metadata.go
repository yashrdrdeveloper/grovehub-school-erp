package metadata

import (
	"context"

	coreevents "grovehub/libs/core/events"
	"grovehub/libs/core/headercontext/internal/src/contracts"
)

type contextKey struct{}

// WithValues attaches trusted metadata and propagates a nonempty event correlation ID.
func WithValues(ctx context.Context, values contracts.Values) context.Context {
	ctx = context.WithValue(ctx, contextKey{}, values)
	if values.CorrelationID != "" {
		ctx = coreevents.WithCorrelationID(ctx, values.CorrelationID)
	}
	return ctx
}

func FromContext(ctx context.Context) contracts.Values {
	if values, ok := ctx.Value(contextKey{}).(contracts.Values); ok {
		return values
	}
	return contracts.Values{}
}

// WithoutIdempotencyKey preserves other metadata while clearing the inherited key.
func WithoutIdempotencyKey(ctx context.Context) context.Context {
	values := FromContext(ctx)
	values.IdempotencyKey = ""
	return WithValues(ctx, values)
}
