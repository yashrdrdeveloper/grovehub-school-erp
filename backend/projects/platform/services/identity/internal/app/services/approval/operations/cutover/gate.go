package cutover

import (
	"context"
	"time"

	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/services/approval"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
)

const RecentAuthenticationWindow = 5 * time.Minute

type Gate struct {
	tenantID string
}

func New(tenantID string) *Gate {
	return &Gate{tenantID: tenantID}
}

func (gate *Gate) Require(ctx context.Context) error {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return err
	}
	if gate == nil || gate.tenantID == "" || gate.tenantID != actor.TenantID {
		return authorization.ErrForbidden
	}
	if actor.AuthenticatedAt.IsZero() || actor.AuthenticatedAt.After(time.Now().UTC()) {
		return principal.ErrRecentAuthenticationRequired
	}
	return principal.Recent(ctx, RecentAuthenticationWindow)
}

type Executor struct {
	gate *Gate
	effect approval.Executor
}

var _ approval.Executor = (*Executor)(nil)

func Wrap(gate *Gate, effect approval.Executor) approval.Executor {
	return &Executor{gate: gate, effect: effect}
}

func (executor *Executor) Execute(ctx context.Context, requestID string, version int64) (entity.ApprovalRequest, error) {
	if err := executor.gate.Require(ctx); err != nil {
		return entity.ApprovalRequest{}, err
	}
	if executor.effect == nil {
		return entity.ApprovalRequest{}, authorization.ErrUnavailable
	}
	return executor.effect.Execute(ctx, requestID, version)
}
