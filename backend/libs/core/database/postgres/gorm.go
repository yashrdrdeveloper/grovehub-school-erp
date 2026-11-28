package postgres

import (
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(gormpostgres.Open(databaseURL), &gorm.Config{})
}
