package registration

import (
	"context"

	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
)

// Service is Identity's registration use-case contract.
// It composes operations/register and operations/resend.
type Service interface {
	Register(ctx context.Context, input accountdto.Create, idempotencyKey string) (accountdto.CreateResult, error)
	Resend(ctx context.Context, accountID string) error
}

// InvitationSender is a narrow delivery-intent interface: it does not own
// registration policy, Foundation Person resolution, or account lifecycle.
type InvitationSender interface {
	Queue(ctx context.Context, account entity.Account, purpose entity.TokenPurpose, recipientName string) error
}
