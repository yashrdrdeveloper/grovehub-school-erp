// Package pipeline composes typed, domain-owned event eligibility and outbox writes.
package pipeline

import (
    "context"
    "errors"

    "grovehub/libs/core/events/outbox"
)

var ErrIncompleteRules = errors.New("event pipeline requires ShouldEmit, Build and Enqueue")

// Rules separates domain eligibility, immutable event creation and durable enqueueing.
type Rules[T any] struct {
    ShouldEmit func(T) bool
    Build func(context.Context, T) (outbox.Message, error)
    Enqueue func(context.Context, outbox.Message) error
}

// Run is called inside the same transaction as the business write.
func Run[T any](ctx context.Context, value T, rules Rules[T]) error {
    if rules.ShouldEmit == nil || rules.Build == nil || rules.Enqueue == nil { return ErrIncompleteRules }
    if !rules.ShouldEmit(value) { return nil }
    message, err := rules.Build(ctx, value)
    if err != nil { return err }
    return rules.Enqueue(ctx, message)
}
