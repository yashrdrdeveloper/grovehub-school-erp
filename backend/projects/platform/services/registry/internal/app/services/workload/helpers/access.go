package helpers

import (
	"context"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/platform/scope"
	definitions "grovehub/services/registry/internal/app/constants/error"
)

// Provisioning is initially platform-operator-only, including cross-project grants.
// Project owners cannot mint or elevate their own workload privileges.
func Require(ctx context.Context, capability string) error {
	if _, err := scope.Require(ctx, scope.Target{Level: scope.Platform}, capability); err != nil {
		return coreerrors.Wrap(definitions.Forbidden, err)
	}
	return nil
}
