package maildelivery

import "context"

type Record struct {
	RequestHash string
	State string
}

const (
	Dispatching = "DISPATCHING"
	Accepted = "ACCEPTED"
)

type Repository interface {
	Find(ctx context.Context, tenantID string, deliveryID string) (Record, bool, error)
	Reserve(ctx context.Context, tenantID string, deliveryID string, requestHash string) (bool, error)
	Accept(ctx context.Context, tenantID string, deliveryID string, requestHash string) error
}
