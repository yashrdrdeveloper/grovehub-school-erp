package services_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	client "grovehub/services/integration/client/connection"
	"grovehub/services/integration/internal/app/repository/maildelivery"
	"grovehub/services/integration/internal/app/services/connection/operations/action/delivery"
)

type deliveryStore struct {
	mu sync.Mutex
	records map[string]maildelivery.Record
}

func (store *deliveryStore) Find(ctx context.Context, tenant string, id string) (maildelivery.Record, bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	record, found := store.records[tenant+"/"+id]
	return record, found, nil
}

func (store *deliveryStore) Reserve(ctx context.Context, tenant string, id string, hash string) (bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := tenant+"/"+id
	if _, exists := store.records[key]; exists { return false, nil }
	store.records[key] = maildelivery.Record{RequestHash: hash, State: maildelivery.Dispatching}
	return true, nil
}

func (store *deliveryStore) Accept(ctx context.Context, tenant string, id string, hash string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.records[tenant+"/"+id] = maildelivery.Record{RequestHash: hash, State: maildelivery.Accepted}
	return nil
}

func TestMailDispatchReplay(t *testing.T) {
	ctx := context.Background()
	store := &deliveryStore{records: map[string]maildelivery.Record{}}
	dispatcher := delivery.New(store)
	var calls atomic.Int32
	send := func() (any,error) { calls.Add(1); return nil,nil }
	if _, err := dispatcher.Dispatch(ctx,"tenant","id",[]byte("message"),send); err != nil { t.Fatal(err) }
	if _, err := dispatcher.Dispatch(ctx,"tenant","id",[]byte("message"),send); err != nil { t.Fatal(err) }
	if _, err := dispatcher.Dispatch(ctx,"tenant","id",[]byte("changed"),send); !errors.Is(err,client.ErrMailDeliveryConflict) { t.Fatalf("expected conflict: %v",err) }
	if _, err := dispatcher.Dispatch(ctx,"other","id",[]byte("message"),send); err != nil { t.Fatal(err) }
	if calls.Load()!=2 { t.Fatal("unexpected provider call count") }
}

func TestMailDispatchUncertain(t *testing.T) {
	ctx := context.Background()
	dispatcher := delivery.New(&deliveryStore{records: map[string]maildelivery.Record{}})
	var calls atomic.Int32
	send := func() (any,error) { calls.Add(1); return nil,errors.New("provider timeout") }
	for range 2 {
		if _, err := dispatcher.Dispatch(ctx,"tenant","id",[]byte("message"),send); !errors.Is(err,client.ErrMailDeliveryUncertain) { t.Fatalf("expected uncertainty: %v",err) }
	}
	if calls.Load()!=1 { t.Fatal("ambiguous dispatch repeated") }
}

func TestMailDispatchConcurrent(t *testing.T) {
	ctx := context.Background()
	dispatcher := delivery.New(&deliveryStore{records: map[string]maildelivery.Record{}})
	started, release := make(chan struct{}), make(chan struct{})
	first := make(chan error,1)
	go func() {
		_, err := dispatcher.Dispatch(ctx,"tenant","id",[]byte("message"),func()(any,error){close(started); <-release; return nil,nil})
		first <- err
	}()
	<-started
	_, err := dispatcher.Dispatch(ctx,"tenant","id",[]byte("message"),func()(any,error){t.Error("duplicate provider invocation"); return nil,nil})
	close(release)
	if err := <-first; err != nil { t.Fatal(err) }
	if !errors.Is(err,client.ErrMailDeliveryUncertain) { t.Fatalf("expected pending reservation: %v",err) }
}
