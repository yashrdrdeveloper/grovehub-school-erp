package unitofwork

import "context"

// UnitOfWork coordinates repository operations within one context-bound transaction.
type UnitOfWork interface {
	Run(ctx context.Context, operation func(context.Context) error) error
}
