package outbox

import (
    "context"
    "encoding/json"
    "errors"
    "log/slog"
    "io"
    "testing"
    "time"
)

type fakeStore struct { entries []Entry; sent []string; retries []string }
func (s *fakeStore) Enqueue(context.Context, Message) error { return nil }
func (s *fakeStore) Claim(context.Context,int,time.Duration)([]Entry,error){ return s.entries,nil }
func (s *fakeStore) MarkSent(_ context.Context,id string) error { s.sent=append(s.sent,id);return nil }
func (s *fakeStore) MarkRetry(_ context.Context,id string,_ int) error { s.retries=append(s.retries,id);return nil }

type fakePublisher struct{ err error; calls int; metadata map[string]string }
func (p *fakePublisher) Publish(_ context.Context,_ string,_ []byte,metadata map[string]string)error{p.calls++;p.metadata=metadata;return p.err}

func TestDispatcherMarksSuccess(t *testing.T) {
    headers,_:=json.Marshal(map[string]string{"event-id":"e1"})
    store:=&fakeStore{entries:[]Entry{{EventID:"e1",Topic:"person",Payload:[]byte("{}"),Metadata:headers,Attempts:1}}}
    transport:=&fakePublisher{}
    dispatcher:=NewDispatcher(store,transport,slog.New(slog.NewTextHandler(io.Discard,nil)))
    if err:=dispatcher.Dispatch(context.Background());err!=nil{t.Fatal(err)}
    if transport.calls!=1||transport.metadata["event-id"]!="e1"||len(store.sent)!=1||len(store.retries)!=0{t.Fatal("delivery was not marked sent")}
}

func TestDispatcherRetriesFailures(t *testing.T) {
    store:=&fakeStore{entries:[]Entry{{EventID:"e2",Topic:"person",Payload:[]byte("{}"),Metadata:[]byte("{}"),Attempts:2}}}
    transport:=&fakePublisher{err:errors.New("broker down")}
    dispatcher:=NewDispatcher(store,transport,slog.New(slog.NewTextHandler(io.Discard,nil)))
    if err:=dispatcher.Dispatch(context.Background());err==nil{t.Fatal("must surface publish errors")}
    if len(store.sent)!=0||len(store.retries)!=1||store.retries[0]!="e2"{t.Fatal("event must remain retryable")}
}
