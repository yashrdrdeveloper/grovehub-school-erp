package strategies

import (
    "context"

    flows "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox/flows"
)

// @Strategy
type ChangedStrategy interface {
    Build(ctx context.Context, payload flows.Payload) (string, []byte, map[string]string, error)
}
