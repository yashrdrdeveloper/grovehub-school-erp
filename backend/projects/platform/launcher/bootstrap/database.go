package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
 "strconv"
 "time"
	"strings"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

const databaseURLEnv = "DATABASE_URL"

var ErrDatabaseURLRequired = errors.New("DATABASE_URL is required")

func NewPostgresDB(ctx context.Context) (*gorm.DB, error) {
	databaseURL, err := databaseURL()
	if err != nil {
		return nil, err
	}

	return openPostgres(ctx,databaseURL,"PLATFORM_DB",20,5)
}

func poolLimits(prefix string, defaultOpen,defaultIdle int)(int,int,error){
 read:=func(suffix string,fallback,min,max int)(int,error){
  raw:=strings.TrimSpace(os.Getenv(prefix+suffix));if raw==""{return fallback,nil}
  n,err:=strconv.Atoi(raw);if err!=nil||n<min||n>max{return 0,errors.New("invalid "+prefix+suffix)}
  return n,nil
 }
 open,err:=read("_MAX_OPEN",defaultOpen,1,100);if err!=nil{return 0,0,err}
 idle,err:=read("_MAX_IDLE",min(defaultIdle,open),0,open);if err!=nil{return 0,0,err}
 return open,idle,nil
}
func openPostgres(ctx context.Context,databaseURL,prefix string,defaultOpen,defaultIdle int) (*gorm.DB,error) {
 open,idle,err:=poolLimits(prefix,defaultOpen,defaultIdle);if err!=nil{return nil,err}
	db, err := coredatabase.NewDB(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres gorm db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("access postgres sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(open)
	sqlDB.SetMaxIdleConns(idle)
	sqlDB.SetConnMaxIdleTime(5*time.Minute)
	sqlDB.SetConnMaxLifetime(30*time.Minute)
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
	if dsn=="" { db,err=NewPostgresDB(ctx) } else { db,err=openPostgres(ctx,dsn,"AUDIT_DB",5,2) }
	if err!=nil { return nil,nil,errors.New("could not open Audit PostgreSQL database") }
	sqlDB,err:=db.DB()
	if err!=nil { return nil,nil,errors.New("could not access Audit PostgreSQL pool") }
	return db,func(){ _=sqlDB.Close() },nil
}
