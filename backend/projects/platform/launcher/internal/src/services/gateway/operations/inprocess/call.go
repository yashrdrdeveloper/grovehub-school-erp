// Package inprocess owns the typed in-process gateway bridge.
// It never constructs identities from untrusted client-supplied data.
package inprocess

import (
	"context"
	"errors"

	schoolgateway "grovehub/libs/platform/session"
)

// Authorizer decides whether a verified actor may invoke an internal operation.
type Authorizer interface {
	Authorize(ctx context.Context, principal schoolgateway.Principal, operation string) error
}

type Operation interface {
	Call(context.Context, string, func(context.Context) error) error
}

// Invoke bridges a typed client operation to a verified in-process caller.
// Generic functions preserve return types without using reflection or any.
func Invoke[T any](ctx context.Context, gateway Operation, name string, invoke func(context.Context) (T, error)) (T, error) {
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
