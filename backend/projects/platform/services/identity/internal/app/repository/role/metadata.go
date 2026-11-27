package role

import "context"

// DraftInspector must be used within the tenant authority transaction.
type DraftInspector interface {
	Require(ctx context.Context, tenantID string, roleID string, version int64) error
}
