package transaction

import (
	"context"

	"gorm.io/gorm"
)

/**
 * UnitOfWork coordinates multiple repository operations in one database
 * transaction while keeping the transaction handle out of service contracts.
 * Nested calls join the transaction already associated with the context.
 */
type UnitOfWork interface {
	Run(ctx context.Context, operation func(context.Context) error) error
}

type unitOfWork struct {
	db *gorm.DB
}

/** NewUnitOfWork creates a reusable GORM-backed transaction coordinator. */
func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &unitOfWork{db: db}
}

func (unit *unitOfWork) Run(ctx context.Context, operation func(context.Context) error) error {
	return Run(ctx, unit.db, operation)
}
