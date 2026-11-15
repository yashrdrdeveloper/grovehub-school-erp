package impl

import (
    "context"
    "errors"
    "testing"

    coreevents "grovehub/libs/core/events"
    coremessaging "grovehub/libs/core/events/messaging"
    messagingdto "grovehub/services/infrastructure/internal/app/dto/messaging"
    "grovehub/services/infrastructure/internal/app/services/messaging/strategies"
)

type stubEventStore struct {
    messages []coreevents.Message
    err error
}

func (store *stubEventStore) Enqueue(_ context.Context, event coreevents.Message) error {
    if store.err != nil { return store.err }
    store.messages=append(store.messages,event)
    return nil
}

type stubStrategy struct {
    published int
    topic string
    payload []byte
    metadata map[string]string
}

func (strategy *stubStrategy) Publish(_ context.Context, topic string, payload []byte, metadata map[string]string) error {
    strategy.published++
    strategy.topic=topic
    strategy.payload=append([]byte(nil),payload...)
    strategy.metadata=metadata
    return nil
}
func (strategy *stubStrategy) Consume(_ context.Context, _ string, _ messagingdto.MessageHandler) error { return nil }
func (strategy *stubStrategy) Health(context.Context) error { return nil }
func (strategy *stubStrategy) Close() error { return nil }

func testMessagingService(store *stubEventStore, broker *stubStrategy) *messagingService {
    registry:=strategies.NewRegistry[strategies.Strategy]()
    registry.Register("test",broker)
    return &messagingService{provider:"test",registry:registry,store:store}
}

func TestPublishWithOptionsDirectDoesNotWriteOutbox(t *testing.T) {
    store:=&stubEventStore{}
    broker:=&stubStrategy{}
    service:=testMessagingService(store,broker)
    if err:=service.PublishWithOptions(context.Background(),"school.test",[]byte("event"),nil,coremessaging.PublishOptions{Persist:false});err!=nil { t.Fatal(err) }
    if broker.published!=1 { t.Fatalf("expected direct NATS publish, got %d",broker.published) }
    if len(store.messages)!=0 { t.Fatal("direct message unexpectedly persisted") }
}

func TestPublishWithOptionsPersistQueuesWithoutPublishing(t *testing.T) {
    store:=&stubEventStore{}
    broker:=&stubStrategy{}
    service:=testMessagingService(store,broker)
    metadata:=map[string]string{"source":"school","event-id":"event-123"}
    payload:=[]byte("original")
    if err:=service.PublishWithOptions(context.Background(),"school.test",payload,metadata,coremessaging.PublishOptions{Persist:true});err!=nil { t.Fatal(err) }
    if broker.published!=0 { t.Fatal("persistent event published synchronously") }
    if len(store.messages)!=1 { t.Fatalf("expected one queued event, got %d",len(store.messages)) }
    got:=store.messages[0]
    if got.EventID!="event-123" || got.Topic!="school.test" || string(got.Payload)!="original" || got.Metadata["source"]!="school" { t.Fatalf("unexpected persisted message: %#v",got) }
    payload[0]='X'
    metadata["source"]="modified"
    if string(got.Payload)!="original" || got.Metadata["source"]!="school" { t.Fatal("persisted event was mutated by caller") }
}

func TestPublishWithOptionsPersistGeneratesEventID(t *testing.T) {
    store:=&stubEventStore{}
    service:=testMessagingService(store,&stubStrategy{})
    if err:=service.PublishWithOptions(context.Background(),"school.test",[]byte("event"),nil,coremessaging.PublishOptions{Persist:true});err!=nil { t.Fatal(err) }
    if len(store.messages)!=1 || store.messages[0].EventID=="" || store.messages[0].Metadata["event-id"]!=store.messages[0].EventID { t.Fatal("missing generated event ID in stored metadata") }
}

func TestPublishWithOptionsPropagatesEnqueueFailure(t *testing.T) {
    want:=errors.New("database unavailable")
    store:=&stubEventStore{err:want}
    broker:=&stubStrategy{}
    service:=testMessagingService(store,broker)
    if err:=service.PublishWithOptions(context.Background(),"school.test",[]byte("event"),nil,coremessaging.PublishOptions{Persist:true});!errors.Is(err,want) { t.Fatalf("wanted database error, got %v",err) }
    if broker.published!=0 { t.Fatal("published even though durable enqueue failed") }
}

func TestPublishRemainsDirect(t *testing.T) {
    store:=&stubEventStore{}
    broker:=&stubStrategy{}
    service:=testMessagingService(store,broker)
    if err:=service.Publish(context.Background(),"school.test",[]byte("event"),nil);err!=nil { t.Fatal(err) }
    if broker.published!=1 || len(store.messages)!=0 { t.Fatal("legacy Publish should remain direct") }
}
