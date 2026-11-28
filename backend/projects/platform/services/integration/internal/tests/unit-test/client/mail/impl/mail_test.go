package impl_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"grovehub/services/integration/client/mail"
	clientimpl "grovehub/services/integration/internal/app/client/mail/impl"
)

type searcher func(context.Context, mail.SearchRequest) (mail.Report, error)

func (s searcher) Search(ctx context.Context, request mail.SearchRequest) (mail.Report, error) {
	return s(ctx, request)
}

func TestMailClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	request := mail.SearchRequest{TenantID: "tenant", After: "cursor", Limit: 10}
	expected := mail.Report{TenantID: "tenant", Rows: []mail.Row{{DeliveryID: "delivery"}}}
	for _, failure := range []error{nil, errors.New("read failed")} {
		client := clientimpl.NewMailClient(searcher(func(got context.Context, input mail.SearchRequest) (mail.Report, error) {
			if got != ctx || input != request {
				t.Fatal("diagnostic scope or context changed")
			}
			return expected, failure
		}))
		report, err := client.Search(ctx, request)
		if !reflect.DeepEqual(report, expected) || err != failure {
			t.Fatalf("report=%v error=%v", report, err)
		}
	}
}
