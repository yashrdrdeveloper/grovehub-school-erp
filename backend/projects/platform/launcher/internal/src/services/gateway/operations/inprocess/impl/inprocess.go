package impl

import (
	"context"
	"errors"

	schoolgateway "grovehub/internal/session"
	"grovehub/launcher/internal/src/services/gateway/operations/inprocess"
)

// @Service
type operation struct {
	authorizer inprocess.Authorizer
}

func New(authorizer inprocess.Authorizer) inprocess.Operation {
	return &operation{authorizer: authorizer}
}

func (gateway *operation) Call(ctx context.Context, name string, invoke func(context.Context) error) error {
	principal, ok := schoolgateway.PrincipalFrom(ctx)
	if !ok {
		return errors.New("trusted gateway principal is required")
	}

	if gateway.authorizer == nil {
		return errors.New("internal operation authorization is not configured")
	}

	if err := gateway.authorizer.Authorize(ctx, principal, name); err != nil {
		return err
	}

	return invoke(ctx)
}

var _ inprocess.Operation = (*operation)(nil)
