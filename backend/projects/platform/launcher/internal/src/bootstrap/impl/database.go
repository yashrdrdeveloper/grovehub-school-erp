package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/launcher/internal/src/config/database/helpers"
)

const databaseURLEnv = "DATABASE_URL"

var ErrDatabaseURLRequired = errors.New("DATABASE_URL is required")

// PostgresPools owns business and Audit pool setup. The opener supplies the
// GORM transport, while limits, health checks and ownership remain here.
type PostgresPools struct {
	open func(string) (*gorm.DB, error)
}

func NewPostgresPools(open func(string) (*gorm.DB, error)) PostgresPools {
	return PostgresPools{open: open}
}

func NewPostgresDB(ctx context.Context) (*gorm.DB, error) {
	return NewPostgresPools(coredatabase.NewDB).Business(ctx)
}

func (p PostgresPools) Business(ctx context.Context) (*gorm.DB, error) {
	db, _, err := p.openBusiness(ctx)
	return db, err
}

func (p PostgresPools) openBusiness(ctx context.Context) (*gorm.DB, *sql.DB, error) {
	databaseURL, err := databaseURL()
	if err != nil {
		return nil, nil, err
	}

	return p.openPostgres(ctx, databaseURL, "PLATFORM_DB", 20, 5)
}

func (p PostgresPools) openPostgres(ctx context.Context, databaseURL, prefix string, defaultOpen, defaultIdle int) (*gorm.DB, *sql.DB, error) {
	open, idle, err := helpers.PoolLimits(prefix, defaultOpen, defaultIdle)
	if err != nil {
		return nil, nil, err
	}

	db, err := p.open(databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("create postgres gorm db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("access postgres sql db: %w", err)
	}

	sqlDB.SetMaxOpenConns(open)
	sqlDB.SetMaxIdleConns(idle)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("ping postgres gorm db: %w", err)
	}

	return db, sqlDB, nil
}

func databaseURL() (string, error) {
	databaseURL := strings.TrimSpace(os.Getenv(databaseURLEnv))
	if databaseURL == "" {
		return "", ErrDatabaseURLRequired
	}

	return databaseURL, nil
}

// NewAuditPostgresDB borrows the business pool by default, or owns a separately
// configured pool. Its cleanup never closes borrowed connections.
func NewAuditPostgresDB(ctx context.Context, shared *gorm.DB) (*gorm.DB, func(), error) {
	return NewPostgresPools(coredatabase.NewDB).Audit(ctx, shared)
}

func (p PostgresPools) Audit(ctx context.Context, shared *gorm.DB) (*gorm.DB, func(), error) {
	dsn := strings.TrimSpace(os.Getenv("AUDIT_DATABASE_URL"))
	if dsn == "" && shared != nil {
		return shared, func() {}, nil
	}

	var db *gorm.DB
	var raw *sql.DB
	var err error
	if dsn == "" {
		db, raw, err = p.openBusiness(ctx)
	} else {
		db, raw, err = p.openPostgres(ctx, dsn, "AUDIT_DB", 5, 2)
	}

	if err != nil {
		return nil, nil, errors.New("could not open Audit PostgreSQL database")
	}

	return db, func() {
		_ = raw.Close()
	}, nil
}
