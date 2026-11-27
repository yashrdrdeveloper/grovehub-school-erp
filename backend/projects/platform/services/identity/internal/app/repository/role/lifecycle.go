package role

import (
	"context"
	"errors"

	"grovehub/services/identity/internal/app/entity"
)

var (
	ErrInvalidTransition = errors.New("invalid role lifecycle transition")
	ErrAssigned          = errors.New("role has account assignments")
	ErrNotReady          = errors.New("role is not ready for activation")
)

type Transition struct {
	TenantID string
	ActorID  string
	RoleID   string
	Version  int64
	Status   entity.RoleStatus
}

// LifecycleWriter is isolated persistence. Only a guarded, approval-bound
// operation may invoke it; it is not a public service or controller.
type LifecycleWriter interface {
	Transition(ctx context.Context, request Transition) (Detail, error)
}
