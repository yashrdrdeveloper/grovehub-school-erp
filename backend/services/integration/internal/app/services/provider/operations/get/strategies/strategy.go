package strategies

import (
    "context"
    dto "grovehub/services/integration/internal/app/dto/provider"
)

// @Strategy
type Strategy interface {
    Execute(ctx context.Context, id string) (dto.ProviderDetails,error)
}
