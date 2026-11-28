package maildelivery

import (
	"context"
	"time"
)

type Receipt struct {
	DeliveryID string
	State string
	ReservedAt *time.Time
	AcceptedAt *time.Time
}

type ReceiptReader interface {
	ListReceiptPage(context.Context, string, string, int) ([]Receipt, error)
}
