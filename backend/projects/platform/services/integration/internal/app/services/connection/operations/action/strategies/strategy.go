package strategies

import (
	"context"
	"encoding/json"

	providerdto "grovehub/services/integration/internal/app/providers/dto"
)

// @Strategy
type Strategy interface {
	ExecuteAction(ctx context.Context, tenantID, connectionID string, action providerdto.Action, input json.RawMessage) (any, error)
	ExecuteSystemMail(ctx context.Context, tenantID string, input json.RawMessage) (any, error)
	ConfigureSystemMailSender(ctx context.Context, tenantID string, connectionID string) error
}
