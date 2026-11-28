package impl

import (
	"context"

	"gorm.io/gorm"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction"

)

type unitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *unitOfWork {
	return &unitOfWork{db: db}
}

func (unit *unitOfWork) Run(ctx context.Context, operation func(context.Context) error) error {
	return transaction.Run(ctx, unit.db, operation)
}
