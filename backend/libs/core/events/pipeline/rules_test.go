package pipeline

import (
    "context"
    "errors"
    "testing"

    "grovehub/libs/core/events/outbox"
)

func TestRulesDoNotEnqueueUnchangedEvents(t *testing.T) {
    calls := 0
    rules := Rules[int]{
        ShouldEmit: func(value int) bool { return value > 0 },
        Build: func(context.Context, int) (outbox.Message, error) { calls++; return outbox.Message{EventID: "e",Topic:"t",Payload:[]byte("{}")},nil },
        Enqueue: func(context.Context, outbox.Message) error { calls++; return nil },
    }
    if err := Run(context.Background(),0,rules); err != nil { t.Fatal(err) }
    if calls != 0 { t.Fatal("unchanged state must not create a message") }
    if err := Run(context.Background(),1,rules); err != nil { t.Fatal(err) }
    if calls != 2 { t.Fatalf("expected build and enqueue, got %d",calls) }
}

func TestRulesPropagateEnqueueErrors(t *testing.T) {
    sentinel := errors.New("outbox unavailable")
    rules := Rules[int]{
        ShouldEmit: func(int)bool{return true},
        Build: func(context.Context,int)(outbox.Message,error){return outbox.Message{EventID:"e"},nil},
        Enqueue: func(context.Context,outbox.Message)error{return sentinel},
    }
    if err := Run(context.Background(),1,rules); !errors.Is(err,sentinel) { t.Fatalf("missing persistence failure: %v",err) }
}
