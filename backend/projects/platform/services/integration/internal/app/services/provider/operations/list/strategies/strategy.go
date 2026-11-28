package strategies

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    "context"
)

// @Strategy
type Strategy interface {
    Execute(ctx context.Context) []providerdto.Descriptor
}
