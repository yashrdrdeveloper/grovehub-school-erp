package principal

import (
	"context"
	"errors"
	"time"

	coregateway "grovehub/libs/core/gateway"
)

type Principal = coregateway.Principal

func With(ctx context.Context, principal Principal) context.Context {
	return coregateway.WithPrincipal(ctx, principal)
}

func From(ctx context.Context) (Principal, bool) {
	return coregateway.PrincipalFrom(ctx)
}

func Require(ctx context.Context) (Principal, error) {
	principal, ok := From(ctx)
	if !ok {
		return Principal{}, errors.New("authentication required")
	}
	return principal, nil
}

func RequireAdmin(ctx context.Context) (Principal, error) {
	principal, err := Require(ctx)
	if err != nil {
		return Principal{}, err
	}
	if !principal.Admin || principal.TenantID == "" {
		return Principal{}, errors.New("admin access forbidden")
	}
	return principal, nil
}

func Recent(ctx context.Context, maxAge time.Duration) error {
	principal, err := Require(ctx)
	if err != nil {
		return err
	}
	if principal.AuthenticatedAt.IsZero() || time.Since(principal.AuthenticatedAt) > maxAge {
		return errors.New("recent authentication required")
	}
	return nil
}
