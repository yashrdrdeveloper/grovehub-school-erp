package strategies

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
)

// @Strategy
type Strategy interface {
    Execute(ctx context.Context, tenantID string, input connectiondto.CreateConnection) (connectiondto.Connection,error)
}
