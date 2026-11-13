package service

import (
    "context"
    "errors"
    "testing"
)

type fakeUnitOfWork struct { calls int; failed bool }
func (unit *fakeUnitOfWork) Run(ctx context.Context, fn func(context.Context) error) error {
    unit.calls++
    err := fn(ctx)
    unit.failed = err != nil
    return err
}

func TestTransactionFailureAndPostRollbackRecovery(t *testing.T) {
    sentinel := errors.New("duplicate")
    unit := &fakeUnitOfWork{}
    state := new(int)
    handler := Transactional[int, string](unit, func(context.Context, *int) Outcome[string] { return Fail[string](sentinel) })
    outcome := Recover(context.Background(), state, handler(context.Background(), state),
        OnError[int, string](sentinel, func(context.Context, *int, error) Outcome[string] {
            if !unit.failed { t.Fatal("recovery executed before rollback") }
            return Return("existing", false)
        }),
    )
    if unit.calls != 1 || outcome.err != nil || outcome.value != "existing" || outcome.created { t.Fatalf("unexpected outcome: %+v", outcome) }
}
