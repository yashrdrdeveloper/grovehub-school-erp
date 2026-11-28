package service

import (
    "context"
    "errors"
)

// Rule describes a predicate and its corresponding typed strategy.
type Rule[S any, R any] struct {
    matches func(*S) bool
    handler Handler[S, R]
}

func When[S any, R any](predicate func(*S) bool, handler Handler[S, R]) Rule[S, R] {
    return Rule[S, R]{matches: predicate, handler: handler}
}

func Otherwise[S any, R any](handler Handler[S, R]) Rule[S, R] {
    return Rule[S, R]{matches: func(*S) bool { return true }, handler: handler}
}

// Decide picks the first matching rule; later rules are not evaluated.
func Decide[S any, R any](ctx context.Context, state *S, rules ...Rule[S, R]) Outcome[R] {
    for _, rule := range rules {
        if rule.matches != nil && rule.matches(state) { return rule.handler(ctx, state) }
    }
    return Fail[R](ErrNoMatchingRule)
}

// Result centralizes Go's error-first branching for a typed operation.
func Result[T any, R any](value T, err error, success func(T) Outcome[R], failure func(error) Outcome[R]) Outcome[R] {
    if err != nil { return failure(err) }
    return success(value)
}

// Lookup distinguishes a missing record from a real storage error.
func Lookup[T any, R any](value T, err error, missing error, found func(T) Outcome[R], notFound func() Outcome[R], failed func(error) Outcome[R]) Outcome[R] {
    if errors.Is(err, missing) { return notFound() }
    return Result(value, err, found, failed)
}
