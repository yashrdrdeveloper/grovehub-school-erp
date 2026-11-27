package protect

import (
	"context"
	"strings"

	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/services/authorization"
)

type Guard struct {
	locker authority.Locker
	reader authority.AdministratorReader
}

var _ authority.Guard = (*Guard)(nil)

func New(locker authority.Locker, reader authority.AdministratorReader) authority.Guard {
	return &Guard{locker: locker, reader: reader}
}

// Run keeps the check and mutation under one tenant authority lock. A failed
// postcondition rolls back the mutation, provided the caller propagates errors.
func (guard *Guard) Run(ctx context.Context, tenantID string, operation func(context.Context) error) error {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID || operation == nil {
		return authority.ErrScopeRequired
	}
	if guard.locker == nil || guard.reader == nil {
		return authorization.ErrUnavailable
	}
	return guard.locker.Run(ctx, tenantID, func(tx context.Context) error {
		before, err := guard.reader.Count(tx, tenantID)
		if err != nil {
			return err
		}
		if before < 1 {
			return authority.ErrBootstrapRequired
		}
		if err := operation(tx); err != nil {
			return err
		}
		after, err := guard.reader.Count(tx, tenantID)
		if err != nil {
			return err
		}
		if after < 1 {
			return authority.ErrLastAdministrator
		}
		return nil
	})
}
