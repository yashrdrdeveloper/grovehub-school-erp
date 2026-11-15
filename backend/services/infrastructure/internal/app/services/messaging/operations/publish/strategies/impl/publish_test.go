package impl

import (
    "context"
    "errors"
    "testing"

    coreevents "grovehub/libs/core/events"
    coremessaging "grovehub/libs/core/events/messaging"
    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
)

type testStore struct {
    events []coreevents.Message
    ctx context.Context
    fail error
}

func (store *testStore) Enqueue(ctx context.Context, event coreevents.Message) error {
    store.ctx = ctx
    if store.fail != nil { return store.fail }
    store.events = append(store.events,event)
    return nil
}

type testProvider struct {
    count int
    topic string
}

func (provider *testProvider) Publish(_ context.Context, topic string, _ []byte, _ map[string]string) error {
    provider.count++
    provider.topic = topic
    return nil
}
func (*testProvider) Consume(context.Context, string, dto.MessageHandler) error { return nil }
func (*testProvider) Health(context.Context) error { return nil }
func (*testProvider) Close() error { return nil }

func TestDirectPublishingBypassesPostgres(t *testing.T) {
    provider := &testProvider{}
    store := &testStore{}
    strategy := New(provider,store)
    ctx := context.Background()
    if err := strategy.Publish(ctx,"foundation_person_changed",[]byte("data"),nil);err!=nil { t.Fatal(err) }
    if err := strategy.PublishWithOptions(ctx,"foundation_person_changed",[]byte("data"),nil,coremessaging.PublishOptions{Persist:false});err!=nil { t.Fatal(err) }
    if provider.count!=2 || len(store.events)!=0 { t.Fatalf("direct publish: count=%d stored=%d",provider.count,len(store.events)) }
}

func TestPersistQueuesFullPayloadWithoutPublishing(t *testing.T) {
    provider := &testProvider{}
    store := &testStore{}
    strategy := New(provider,store)
    ctx := context.WithValue(context.Background(),struct{}{},"transaction-context")
    metadata := map[string]string{"event-id":"event-123","source":"foundation"}
    payload := []byte("complete old/new entity snapshots")
    if err := strategy.PublishWithOptions(ctx,"foundation_person_changed",payload,metadata,coremessaging.PublishOptions{Persist:true});err!=nil { t.Fatal(err) }
    if provider.count!=0 || len(store.events)!=1 { t.Fatalf("persistent publish: count=%d stored=%d",provider.count,len(store.events)) }
    got:=store.events[0]
    if store.ctx!=ctx || got.EventID!="event-123" || got.Topic!="foundation_person_changed" || string(got.Payload)!=string(payload) { t.Fatalf("bad persisted event: %#v",got) }
    payload[0]='X'
    metadata["source"]="mutated"
    if got.Metadata["source"]!="foundation" || string(got.Payload)!="complete old/new entity snapshots" { t.Fatal("event was mutated by caller") }
}

func TestPersistGeneratesStableEventID(t *testing.T) {
    store := &testStore{}
    strategy := New(&testProvider{},store)
    if err:=strategy.PublishWithOptions(context.Background(),"topic",[]byte("snapshot"),nil,coremessaging.PublishOptions{Persist:true});err!=nil { t.Fatal(err) }
    if len(store.events)!=1 || store.events[0].EventID=="" || store.events[0].Metadata["event-id"]!=store.events[0].EventID { t.Fatal("event ID missing") }
}

func TestPersistReturnsStoreErrorWithoutPublishing(t *testing.T) {
    provider := &testProvider{}
    want := errors.New("database unavailable")
    strategy := New(provider,&testStore{fail:want})
    err:=strategy.PublishWithOptions(context.Background(),"topic",[]byte("snapshot"),nil,coremessaging.PublishOptions{Persist:true})
    if !errors.Is(err,want) || provider.count!=0 { t.Fatalf("want enqueue error without publish, got %v, count %d",err,provider.count) }
}
