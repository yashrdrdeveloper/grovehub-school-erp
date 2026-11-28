package invitation

import "time"

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
