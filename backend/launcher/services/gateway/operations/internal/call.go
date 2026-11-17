// Package internal is the in-process, typed gateway bridge.
// It does not convert untrusted HTTP headers into service identities.
package internal

import (
	"context"
	"errors"

	coregateway "grovehub/libs/core/gateway"
)

type Authorizer interface {
	Authorize(ctx context.Context, principal coregateway.Principal, operation string) error
}

type Operation struct {
	authorizer Authorizer
}

func New(authorizer Authorizer) *Operation {
	return &Operation{authorizer: authorizer}
}

func (gateway *Operation) Call(ctx context.Context, name string, invoke func(context.Context) error) error {
	principal, ok := coregateway.PrincipalFrom(ctx)
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
