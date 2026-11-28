package contracts

import (
	"context"
	"errors"
)

var ErrIncompleteRules = errors.New("event pipeline requires ShouldEmit, Build and Enqueue")

// Rules composes eligibility, construction and delivery independently of transport.
// M is the consumer-owned message type.
type Rules[T any, M any] struct {
	ShouldEmit func(T) bool
	Build      func(context.Context, T) (M, error)
	Enqueue    func(context.Context, M) error
}
