package authorization

import (
	"context"

	"grovehub/libs/platform/scope"
	"grovehub/services/registry/client/authorization"
)

type AuthorizationRepository interface {
	Resolve(ctx context.Context, keyID string, target scope.Target, capability string, allowInactiveTarget bool) (authorization.Credential, error)
	TenantState(ctx context.Context, projectID, tenantID string) (authorization.TenantState, error)
}
