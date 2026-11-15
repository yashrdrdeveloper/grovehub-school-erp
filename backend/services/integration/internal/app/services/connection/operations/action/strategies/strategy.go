package strategies

import (
    "context"
    "encoding/json"
    "grovehub/services/integration/internal/app/providers"
)

// @Strategy
type Strategy interface {
    ExecuteAction(ctx context.Context, tenantID, connectionID string, action providers.Action, input json.RawMessage) (any,error)
}
