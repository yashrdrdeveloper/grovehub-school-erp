package impl

import (
	"gorm.io/gorm"

	"grovehub/services/audit/internal/app/repository/schema"
	"grovehub/services/audit/internal/app/repository/schema/constants"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) schema.Repository {
	return &repository{db: db}
}

func (r *repository) CheckActivity() error {
	return r.db.Exec(constants.ActivityQuery).Error
}

func (r *repository) CheckRejection() error {
	return r.db.Exec(constants.RejectionQuery).Error
}

func (r *repository) CheckReceipt() error {
	return r.db.Exec(constants.ReceiptQuery).Error
}

func (r *repository) CheckOperator() error {
	return r.db.Exec(constants.OperatorQuery).Error
}

func (r *repository) CheckSchoolReceipt() error {
	return r.db.Exec(constants.SchoolReceiptQuery).Error
}
