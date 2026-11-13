package service

import (
    "context"
    "errors"

    coretransaction "grovehub/libs/core/database/transaction"
)

var ErrMissingUnitOfWork = errors.New("service transaction requires a unit of work")

type Operation func(context.Context) error

// Sequence runs dependent persistence operations and stops at the first error.
func Sequence[R any](ctx context.Context, result R, created bool, operations ...Operation) Outcome[R] {
    for _, operation := range operations {
        if err := operation(ctx); err != nil { return Fail[R](err) }
    }
    return Return(result, created)
}

// Transactional ensures all nested repository writes share a transaction.
// The annotation comment alone is not executable; callers must use this wrapper.
func Transactional[S any, R any](unit coretransaction.UnitOfWork, handler Handler[S, R]) Handler[S, R] {
    return func(ctx context.Context, state *S) Outcome[R] {
        if unit == nil { return Fail[R](ErrMissingUnitOfWork) }
        var result Outcome[R]
        err := unit.Run(ctx, func(txCtx context.Context) error {
            result = handler(txCtx, state)
            return result.err
        })
        if err != nil { return Fail[R](err) }
        return result
    }
}

// Recovery is evaluated only after the transaction wrapper has returned.
// This avoids querying using the context of an aborted Postgres transaction.
type Recovery[S any, R any] struct {
    matches func(error) bool
    handler func(context.Context, *S, error) Outcome[R]
}

func OnError[S any, R any](target error, handler func(context.Context, *S, error) Outcome[R]) Recovery[S, R] {
    return Recovery[S, R]{matches: func(err error) bool { return errors.Is(err, target) }, handler: handler}
}

func OtherwiseError[S any, R any](handler func(context.Context, *S, error) Outcome[R]) Recovery[S, R] {
    return Recovery[S, R]{matches: func(error) bool { return true }, handler: handler}
}

func Recover[S any, R any](ctx context.Context, state *S, result Outcome[R], recoveries ...Recovery[S, R]) Outcome[R] {
    if result.action != failChain { return result }
    for _, recovery := range recoveries {
        if recovery.matches != nil && recovery.matches(result.err) {
            return recovery.handler(ctx, state, result.err)
        }
    }
    return result
}
