package helpers

import (
	"database/sql"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/constants"
)

func ConfigurePool(db *sql.DB) {
	db.SetMaxOpenConns(constants.MaxOpenConns)
	db.SetMaxIdleConns(constants.MaxIdleConns)
	db.SetConnMaxIdleTime(constants.ConnMaxIdleTime)
	db.SetConnMaxLifetime(constants.ConnMaxLifetime)
}
