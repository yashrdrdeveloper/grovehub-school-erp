// Package pipeline composes transport-agnostic, typed event eligibility and delivery rules.
package pipeline

import (
    "context"
    "errors"
)

var ErrIncompleteRules = errors.New("event pipeline requires ShouldEmit, Build and Enqueue")

// M keeps core independent from Infrastructure-owned messaging DTOs.
type Rules[T any, M any] struct {
    ShouldEmit func(T) bool
    Build func(context.Context, T) (M, error)
    Enqueue func(context.Context, M) error
}

func Run[T any, M any](ctx context.Context, value T, rules Rules[T,M]) error {
    if rules.ShouldEmit==nil || rules.Build==nil || rules.Enqueue==nil {return ErrIncompleteRules}
    if !rules.ShouldEmit(value){return nil}
    message,err:=rules.Build(ctx,value)
    if err!=nil{return err}
    return rules.Enqueue(ctx,message)
}
