package registration

import (
	"context"

	"grovehub/services/identity/internal/app/entity"
)

// Service owns invitation issuance and retry-safe delivery intent.
// Login, account lifecycle, and credential hashing are separate concerns.
type Service interface {
	Queue(ctx context.Context, account entity.Account, purpose entity.TokenPurpose, recipientName string) error
}
