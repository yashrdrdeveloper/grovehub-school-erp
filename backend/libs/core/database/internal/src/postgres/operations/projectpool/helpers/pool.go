package helpers

import (
	"database/sql"

	"grovehub/libs/core/database/internal/src/postgres/operations/projectpool/constants"
)

func ConfigurePool(db *sql.DB) {
	db.SetMaxOpenConns(constants.MaxOpenConns)
	db.SetMaxIdleConns(constants.MaxIdleConns)
	db.SetConnMaxIdleTime(constants.ConnMaxIdleTime)
	db.SetConnMaxLifetime(constants.ConnMaxLifetime)
}
