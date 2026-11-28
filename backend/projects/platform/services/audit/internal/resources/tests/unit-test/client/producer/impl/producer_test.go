package impl_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"grovehub/services/audit/client/producer"
	clientimpl "grovehub/services/audit/internal/app/client/producer/impl"
)

type service struct {
	ctx context.Context
	event producer.Event
	activity producer.Activity
	err error
}

func (s *service) Produce(ctx context.Context, event producer.Event) error {
	s.ctx, s.event = ctx, event
	return s.err
}

func (s *service) Record(ctx context.Context, activity producer.Activity) error {
	s.ctx, s.activity = ctx, activity
	return s.err
}

func TestProducerClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	failure := errors.New("publication failed")
	event := producer.Event{Type: producer.ActivityOccurred, Payload: []byte("payload")}
	activity := producer.Activity{TenantID: "tenant", AssetID: "asset"}
	for _, err := range []error{nil, failure} {
		s := &service{err: err}
		client := clientimpl.NewProducerClient(s)
		if got := client.Produce(ctx, event); got != err || s.ctx != ctx || !reflect.DeepEqual(s.event, event) {
			t.Fatalf("event delegation failed: %v", got)
		}
		if got := client.Record(ctx, activity); got != err || s.ctx != ctx || !reflect.DeepEqual(s.activity, activity) {
			t.Fatalf("activity delegation failed: %v", got)
		}
	}
}
