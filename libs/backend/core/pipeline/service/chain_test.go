package service

import (
    "context"
    "errors"
    "testing"
)

func TestChainStopsAtFirstReturn(t *testing.T) {
    calls := 0
    one := func(context.Context, *int) Outcome[string] { calls++; return Next[string]() }
    two := func(context.Context, *int) Outcome[string] { calls++; return Return("existing", false) }
    three := func(context.Context, *int) Outcome[string] { calls++; return Return("unexpected", true) }
    value, created, err := Run(context.Background(), new(int), one, two, three)
    if err != nil || created || value != "existing" || calls != 2 { t.Fatalf("unexpected chain result: %q, %v, %v (%d)", value, created, err, calls) }
}

func TestChainFailsOnUnresolvedNext(t *testing.T) {
    _, _, err := Run(context.Background(), new(int), Continue[int, string])
    if !errors.Is(err, ErrUnresolvedChain) { t.Fatalf("expected unresolved chain, got %v", err) }
}

func TestChainPropagatesFailure(t *testing.T) {
    sentinel := errors.New("stopped")
    _, _, err := Run(context.Background(), new(int), func(context.Context, *int) Outcome[string] { return Fail[string](sentinel) })
    if !errors.Is(err, sentinel) { t.Fatalf("expected failure, got %v", err) }
}
