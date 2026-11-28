package contracts

import "context"

// UnitOfWork coordinates repository operations in a context-bound transaction.
type UnitOfWork interface {
	Run(ctx context.Context, operation func(context.Context) error) error
}
