package impl

import (
	"context"

	"gorm.io/gorm"

	"grovehub/libs/core/database/internal/src/transaction/operations/run"
	"grovehub/libs/core/database/internal/src/transaction/operations/unitofwork"
)

type unitOfWork struct {
	db *gorm.DB
}

var _ unitofwork.UnitOfWork = (*unitOfWork)(nil)

func NewUnitOfWork(db *gorm.DB) unitofwork.UnitOfWork {
	return &unitOfWork{db: db}
}

func (unit *unitOfWork) Run(ctx context.Context, operation func(context.Context) error) error {
	return run.Run(ctx, unit.db, operation)
}
