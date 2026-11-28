package postgres

import (
	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

func NewDB(databaseURL string) (*gorm.DB, error) {
	return coredatabase.NewDB(databaseURL)
}
