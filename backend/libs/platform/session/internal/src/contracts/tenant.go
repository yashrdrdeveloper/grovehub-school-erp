package contracts

import (
	"context"
	"errors"
)

type TenantCheck func(context.Context, string) (int64, error)

// TenantStateLookup returns authoritative state for a known school.
type TenantState struct {
	Generation int64
	Active     bool
}

type TenantStateLookup func(context.Context, string) (TenantState, error)

var ErrTenantDenied = errors.New("tenant access denied")
