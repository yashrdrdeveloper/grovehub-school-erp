package strategies

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
)

// @Strategy
type Strategy interface {
    RunTest(ctx context.Context, tenantID, id, direction string) (connectiondto.ConnectionTest,error)
    ListTests(ctx context.Context, tenantID, id string) ([]connectiondto.ConnectionTest,error)
}
