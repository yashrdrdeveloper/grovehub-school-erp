package rejection

import "context"

type Repository interface {
	Record(ctx context.Context, eventType string, payload []byte, reason string) error
}
