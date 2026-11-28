package transaction

import (
	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

type UnitOfWork = coredatabase.UnitOfWork

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return coredatabase.NewUnitOfWork(db)
}
