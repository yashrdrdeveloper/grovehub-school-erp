package setup

import (
	"context"
	accountdto "grovehub/services/identity/internal/app/dto/account"
)

// Operation is the stable use-case contract; implementation lives in impl/.
type Operation interface {
	Execute(ctx context.Context, grant string, input accountdto.Setup) (accountdto.SetupResult, error)
}
