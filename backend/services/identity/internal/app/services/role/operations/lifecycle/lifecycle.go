// Package lifecycle provides the guarded internal role transition primitive.
// Approval, delegation revalidation and outbox must be enforced by its future
// caller before any live use. It is not exposed by the role service or routes.
package lifecycle

import (
	"context"

	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/repository/role/helpers"
	"grovehub/services/identity/internal/app/services/authorization"
)

type Operation struct {
	guard  authority.Guard
	writer role.LifecycleWriter
}

func New(guard authority.Guard, writer role.LifecycleWriter) *Operation {
	return &Operation{guard: guard, writer: writer}
}

func (operation *Operation) Execute(ctx context.Context, request role.Transition) (role.Detail, error) {
	request, err := helpers.NormalizeTransition(request)
	if err != nil {
		return role.Detail{}, err
	}
	if operation.guard == nil || operation.writer == nil {
		return role.Detail{}, authorization.ErrUnavailable
	}
	var detail role.Detail
	err = operation.guard.Run(ctx, request.TenantID, func(tx context.Context) error {
		var err error
		detail, err = operation.writer.Transition(tx, request)
		return err
	})
	if err != nil {
		return role.Detail{}, err
	}
	return detail, nil
}
