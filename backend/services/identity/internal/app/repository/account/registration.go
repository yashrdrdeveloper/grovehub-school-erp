package account

import (
	"context"
	"time"
)

// Registration inserts a pending, unassigned account; it has no role or
// activation fields. No invitation or credentials may be created here.
type Registration struct {
	TenantID string
	ActorID string
	AccountID string
	PersonGID string
	Email string
	Phone *string
	Key string
	Digest string
	EventID string
	At time.Time
}

type RegistrationReceipt struct {
	AccountID string
	PersonGID string
	Email string
	Digest string
}

type RegistrationWriter interface {
	FindRegistration(ctx context.Context, tenantID string, actorID string, key string) (RegistrationReceipt, error)
	CreatePending(ctx context.Context, registration Registration) error
}
