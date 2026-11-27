package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"

	corepostgres "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
)

const databaseURLEnv = "DATABASE_URL"

var ErrDatabaseURLRequired = errors.New("DATABASE_URL is required")

func NewPostgresDB(ctx context.Context) (*gorm.DB, error) {
	databaseURL, err := databaseURL()
	if err != nil {
		return nil, err
	}

	return openPostgres(ctx,databaseURL)
}

func openPostgres(ctx context.Context,databaseURL string) (*gorm.DB,error) {
	db, err := corepostgres.NewDB(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres gorm db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("access postgres sql db: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping postgres gorm db: %w", err)
	}

	return db, nil
}

func databaseURL() (string, error) {
	databaseURL := strings.TrimSpace(os.Getenv(databaseURLEnv))
	if databaseURL == "" {
		return "", ErrDatabaseURLRequired
	}
	return databaseURL, nil
}

// NewAuditPostgresDB reuses the business pool by default. Explicit Audit DSN
// opens an independently owned pool, even if it names the same server.
// The returned cleanup never closes the borrowed business pool.
func NewAuditPostgresDB(ctx context.Context,shared *gorm.DB) (*gorm.DB,func(),error) {
	dsn:=strings.TrimSpace(os.Getenv("AUDIT_DATABASE_URL"))
	if dsn=="" && shared!=nil { return shared,func(){},nil }
	var db *gorm.DB
	var err error
	if dsn=="" { db,err=NewPostgresDB(ctx) } else { db,err=openPostgres(ctx,dsn) }
	if err!=nil { return nil,nil,errors.New("could not open Audit PostgreSQL database") }
	sqlDB,err:=db.DB()
	if err!=nil { return nil,nil,errors.New("could not access Audit PostgreSQL pool") }
	return db,func(){ _=sqlDB.Close() },nil
}
