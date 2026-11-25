package principal

import (
	"context"
	"errors"
	"time"

	coregateway "grovehub/libs/core/gateway"
)

type Principal = coregateway.Principal

var (
	ErrAuthenticationRequired = errors.New("authentication required")
	ErrAdminForbidden = errors.New("admin access forbidden")
	ErrTenantForbidden = errors.New("tenant access forbidden")
	ErrRecentAuthenticationRequired = errors.New("recent authentication required")
)

func With(ctx context.Context, principal Principal) context.Context {
	return coregateway.WithPrincipal(ctx, principal)
}

func From(ctx context.Context) (Principal, bool) {
	return coregateway.PrincipalFrom(ctx)
}

func Require(ctx context.Context) (Principal, error) {
	principal, ok := From(ctx)
	if !ok {
		return Principal{}, ErrAuthenticationRequired
	}
	return principal, nil
}

// RequireTenant enforces authenticated tenant membership without a role check.
// Account authorization is deferred to the future RBAC policy.
func RequireTenant(ctx context.Context) (Principal, error) {
	actor, err := Require(ctx)
	if err != nil {
		return Principal{}, err
	}
	if actor.TenantID == "" {
		return Principal{}, ErrTenantForbidden
	}
	return actor, nil
}

func RequireAdmin(ctx context.Context) (Principal, error) {
	principal, err := Require(ctx)
	if err != nil {
		return Principal{}, err
	}
	if !principal.Admin || principal.TenantID == "" {
		return Principal{}, ErrAdminForbidden
	}
	return principal, nil
}

func Recent(ctx context.Context, maxAge time.Duration) error {
	principal, err := Require(ctx)
	if err != nil {
		return err
	}
	if principal.AuthenticatedAt.IsZero() || time.Since(principal.AuthenticatedAt) > maxAge {
		return ErrRecentAuthenticationRequired
	}
	return nil
}
