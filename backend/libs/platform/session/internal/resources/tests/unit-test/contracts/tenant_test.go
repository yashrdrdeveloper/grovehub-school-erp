package contracts_test

import (
	"context"
	"errors"
	"testing"

	"grovehub/libs/platform/session/internal/src/contracts"
)

func TestTenantContracts(t *testing.T) {
	var lookup contracts.TenantStateLookup = func(context.Context, string) (contracts.TenantState, error) {
		return contracts.TenantState{Generation: 7, Active: true}, nil
	}
	var check contracts.TenantCheck = func(ctx context.Context, id string) (int64, error) {
		if id == "" {
			return 0, contracts.ErrTenantDenied
		}
		state, err := lookup(ctx, id)
		return state.Generation, err
	}
	if n, err := check(context.Background(), "school"); n != 7 || err != nil {
		t.Fatalf("tenant state contract: %d, %v", n, err)
	}
	if _, err := check(context.Background(), ""); !errors.Is(err, contracts.ErrTenantDenied) {
		t.Fatalf("denial contract: %v", err)
	}
}
