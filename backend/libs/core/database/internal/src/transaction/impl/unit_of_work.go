package impl

import (
	"context"

	"gorm.io/gorm"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/contracts"
)

type unitOfWork struct {
	db *gorm.DB
}

var _ contracts.UnitOfWork = (*unitOfWork)(nil)

func NewUnitOfWork(db *gorm.DB) contracts.UnitOfWork {
	return &unitOfWork{db: db}
}

func (unit *unitOfWork) Run(ctx context.Context, operation func(context.Context) error) error {
	return transaction.Run(ctx, unit.db, operation)
}
