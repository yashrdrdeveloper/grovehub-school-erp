package registration

import (
	"context"

	repo "grovehub/services/identity/internal/app/repository/account"
)

type Producer interface {
	Enqueue(ctx context.Context, registration repo.Registration) error
}
