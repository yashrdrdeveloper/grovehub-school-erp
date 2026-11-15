package strategies

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    "context"
    "encoding/json"
)

// @Strategy
type Strategy interface {
    ExecuteAction(ctx context.Context, tenantID, connectionID string, action providerdto.Action, input json.RawMessage) (any,error)
}
