// Package service runs typed, ordered service handlers. Only the core runner
// owns branching mechanics; individual services retain their domain policies.
package service

import (
    "context"
    "errors"
)

var ErrUnresolvedChain = errors.New("service chain ended without a result")
var ErrNoMatchingRule = errors.New("service strategy has no matching rule")

type action uint8

const (
    continueChain action = iota
    stopChain
    failChain
)

// Outcome is an explicit Next, Return, or Fail decision from a handler.
type Outcome[R any] struct {
    action action
    value R
    created bool
    err error
}

type Handler[S any, R any] func(context.Context, *S) Outcome[R]

// Strategy allows reusable, stateful handlers without reflection.
type Strategy[S any, R any] interface {
    Execute(context.Context, *S) Outcome[R]
}

func Next[R any]() Outcome[R] { return Outcome[R]{action: continueChain} }
func Return[R any](value R, created bool) Outcome[R] { return Outcome[R]{action: stopChain, value: value, created: created} }
func Fail[R any](err error) Outcome[R] { return Outcome[R]{action: failChain, err: err} }

// Continue adapts Next to a handler for declarative rules.
func Continue[S any, R any](context.Context, *S) Outcome[R] { return Next[R]() }

// Chain stops at the first Return/Fail. Ending with Next is a programming error.
func Chain[S any, R any](ctx context.Context, state *S, handlers ...Handler[S, R]) Outcome[R] {
    for _, handler := range handlers {
        if handler == nil { return Fail[R](ErrUnresolvedChain) }
        outcome := handler(ctx, state)
        if outcome.action != continueChain { return outcome }
    }
    return Fail[R](ErrUnresolvedChain)
}

func Run[S any, R any](ctx context.Context, state *S, handlers ...Handler[S, R]) (R, bool, error) {
    result := Chain(ctx, state, handlers...)
    return result.value, result.created, result.err
}

// RequireHandled converts an unexpected Next during conflict recovery to Fail.
func RequireHandled[R any](outcome Outcome[R], err error) Outcome[R] {
    if outcome.action == continueChain { return Fail[R](err) }
    return outcome
}
