package impl_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"grovehub/services/identity/client/invitation"
	clientimpl "grovehub/services/identity/internal/app/client/invitation/impl"
)

type reader struct {
	check func(context.Context, string, string, int) ([]invitation.Delivery, error)
}

func (r reader) ListDeliveryPage(ctx context.Context, tenantID, after string, limit int) ([]invitation.Delivery, error) {
	return r.check(ctx, tenantID, after, limit)
}

func TestInvitationClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	expected := []invitation.Delivery{{DeliveryID: "delivery", AccountID: "account"}}
	for _, failure := range []error{nil, errors.New("read failed")} {
		client := clientimpl.NewInvitationClient(reader{check: func(got context.Context, tenant, after string, limit int) ([]invitation.Delivery, error) {
			if got != ctx || tenant != "tenant" || after != "cursor" || limit != 10 {
				t.Fatal("invitation scope or context changed")
			}
			return expected, failure
		}})
		rows, err := client.ListDeliveryPage(ctx, "tenant", "cursor", 10)
		if !reflect.DeepEqual(rows, expected) || err != failure {
			t.Fatalf("rows=%v error=%v", rows, err)
		}
	}
}
