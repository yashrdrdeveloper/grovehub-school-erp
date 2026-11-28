// Package invitation exposes the read-only Identity delivery projection.
package invitation

import (
	"context"
	"time"
)

// Delivery intentionally excludes addresses, bodies, tokens and request hashes.
type Delivery struct {
	DeliveryID string
	AccountID string
	Purpose string
	SentAt *time.Time
	SuppressedAt *time.Time
	SuppressionReason *string
	LeaseUntil *time.Time
}

type DeliveryReader interface {
	ListDeliveryPage(ctx context.Context, tenantID, after string, limit int) ([]Delivery, error)
}
