package correlation

import (
	"context"

	"grovehub/libs/core/telemetry"
)

type key struct{}

func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key{}, id)
}

// Resolve prefers explicit metadata, then context, then a valid trace ID.
func Resolve(ctx context.Context, explicit string) string {
	if explicit != "" {
		return explicit
	}
	if value, ok := ctx.Value(key{}).(string); ok && value != "" {
		return value
	}
	return telemetry.TraceID(ctx)
}
