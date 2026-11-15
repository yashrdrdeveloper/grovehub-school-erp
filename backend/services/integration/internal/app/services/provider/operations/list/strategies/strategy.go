package strategies

import (
    "context"
    "grovehub/services/integration/internal/app/providers"
)

// @Strategy
type Strategy interface {
    Execute(ctx context.Context) []providers.Descriptor
}
