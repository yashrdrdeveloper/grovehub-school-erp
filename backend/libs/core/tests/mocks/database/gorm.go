// Package database provides an isolated PostgreSQL-shaped SQL connection
// using sqlmock. No unit test can reach a live database through this helper.
package database

import (
    "database/sql"

    "github.com/DATA-DOG/go-sqlmock"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// NewPostgres builds a GORM connection backed entirely by a fake SQL driver.
// All callers must check ExpectationsWereMet and invoke cleanup.
func NewPostgres() (*gorm.DB, sqlmock.Sqlmock, func() error, error) {
    return NewPostgresWith(func() (*sql.DB, sqlmock.Sqlmock, error) {
        return sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
    }, openPostgres)
}

// NewPostgresWith lets isolated tests supply their own fake driver and GORM opener.
// Production test fixtures use NewPostgres; no actual database is opened.
func NewPostgresWith(connect func() (*sql.DB, sqlmock.Sqlmock, error), open func(*sql.DB) (*gorm.DB, error)) (*gorm.DB, sqlmock.Sqlmock, func() error, error) {
    connection, mock, err := connect()
    if err != nil { return nil,nil,nil,err }
    db,err := open(connection)
    if err != nil { _=connection.Close();return nil,nil,nil,err }
    return db,mock,connection.Close,nil
}

func openPostgres(connection *sql.DB) (*gorm.DB,error) {
    return gorm.Open(postgres.New(postgres.Config{Conn:connection,PreferSimpleProtocol:true}),
        &gorm.Config{SkipDefaultTransaction:true,DisableAutomaticPing:true})
}
