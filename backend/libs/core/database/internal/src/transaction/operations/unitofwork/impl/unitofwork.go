package impl

import (
	"context"

	"gorm.io/gorm"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/run"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/unitofwork"
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
