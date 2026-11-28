package causation

import "context"

type key struct{}

func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key{}, id)
}

// Resolve prefers explicit metadata, then the originating command/event ID.
func Resolve(ctx context.Context, explicit string) string {
	if explicit != "" {
		return explicit
	}
	if value, ok := ctx.Value(key{}).(string); ok {
		return value
	}
	return ""
}
