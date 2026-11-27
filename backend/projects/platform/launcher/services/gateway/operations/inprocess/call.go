// Package inprocess owns the typed in-process gateway bridge.
// It never constructs identities from untrusted client-supplied data.
package inprocess

import (
	"context"
	"errors"

	schoolgateway "grovehub/libs/schoolgateway"
)

// Authorizer decides whether a verified actor may invoke an internal operation.
type Authorizer interface {
	Authorize(ctx context.Context, principal schoolgateway.Principal, operation string) error
}

type Operation struct {
	authorizer Authorizer
}

func New(authorizer Authorizer) *Operation {
	return &Operation{authorizer: authorizer}
}

func (gateway *Operation) Call(ctx context.Context, name string, invoke func(context.Context) error) error {
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

// Invoke bridges a typed client operation to a verified in-process caller.
// Generic functions preserve return types without using reflection or any.
func Invoke[T any](ctx context.Context, gateway *Operation, name string, invoke func(context.Context) (T, error)) (T, error) {
	var result T
	if gateway == nil {
		return result, errors.New("internal gateway is not configured")
	}
	err := gateway.Call(ctx, name, func(authorized context.Context) error {
		var callErr error
		result, callErr = invoke(authorized)
		return callErr
	})
	return result, err
}
