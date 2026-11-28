package maildelivery_test

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"
	"grovehub/operations/mailreconcile"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	client "grovehub/services/integration/client/connection"
	repository "grovehub/services/integration/internal/app/repository/maildelivery/impl"
	"grovehub/services/integration/internal/app/services/connection/operations/action/delivery"
	"grovehub/services/integration/internal/resources/migrations"
)

func TestDurableDispatch(t *testing.T) {
	dsn := os.Getenv("INTEGRATION_MAIL_TEST_DSN")
	if dsn == "" { t.Skip("requires disposable PostgreSQL") }
	ctx := context.Background()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { t.Fatal(err) }
	sqlDB, err := db.DB()
	if err != nil { t.Fatal(err) }
	defer sqlDB.Close()
	var name string
	if err := db.Raw("SELECT current_database()").Scan(&name).Error; err != nil || name != "groovehub_rbac_test" {
		t.Fatal("mail integration test requires disposable RBAC database")
	}
	var existing int64
	if err := db.Raw("SELECT count(*) FROM information_schema.schemata WHERE schema_name = 'integration'").Scan(&existing).Error; err != nil || existing != 0 {
		t.Fatal("mail integration schema must be fresh")
	}
	if err := migrations.Migrate(ctx, sqlDB); err != nil { t.Fatal(err) }
	dispatcher := delivery.New(repository.New(db))
	tenant, id := uuid.NewString(), uuid.NewString()
	var calls atomic.Int32
	send := func() (any,error) { calls.Add(1); return nil,nil }
	rollback := errors.New("injected caller rollback")
	err = coredatabase.RunTransaction(ctx, db, func(tx context.Context) error {
		if _, err := dispatcher.Dispatch(tx, tenant, id, []byte("synthetic mail"), send); err != nil { return err }
		return rollback
	})
	if !errors.Is(err, rollback) { t.Fatalf("expected rollback: %v", err) }
	// A new dispatcher simulates worker restart after Identity rollback.
	dispatcher = delivery.New(repository.New(db))
	if _, err := dispatcher.Dispatch(ctx, tenant, id, []byte("synthetic mail"), send); err != nil { t.Fatal(err) }
	if calls.Load() != 1 { t.Fatal("accepted delivery was repeated") }
	if _, err := dispatcher.Dispatch(ctx, tenant, id, []byte("different mail"), send); !errors.Is(err, client.ErrMailDeliveryConflict) { t.Fatalf("expected conflict: %v",err) }

	unknown := uuid.NewString()
	if _, err := dispatcher.Dispatch(ctx, tenant, unknown, []byte("synthetic mail"), func()(any,error){calls.Add(1); return nil,errors.New("timeout")}); !errors.Is(err,client.ErrMailDeliveryUncertain) { t.Fatalf("expected uncertainty: %v",err) }
	if _, err := dispatcher.Dispatch(ctx, tenant, unknown, []byte("synthetic mail"), send); !errors.Is(err,client.ErrMailDeliveryUncertain) { t.Fatalf("expected persistent uncertainty: %v",err) }
	if calls.Load() != 2 { t.Fatal("uncertain delivery was repeated") }

	// No matching Identity jobs exist for these synthetic Integration receipts.
	report, err := mailreconcile.Read(ctx, db, tenant, "", 1)
	if err != nil { t.Fatal(err) }
	if len(report.Rows) != 1 || report.Rows[0].Disposition != "RECEIPT_WITHOUT_JOB" || report.Rows[0].AccountID != nil || report.NextCursor == "" { t.Fatal("unexpected first reconciliation page") }
	next, err := mailreconcile.Read(ctx, db, tenant, report.NextCursor, 1)
	if err != nil || len(next.Rows) != 1 || next.Rows[0].DeliveryID == report.Rows[0].DeliveryID || next.NextCursor != "" { t.Fatalf("invalid continuation page: %v", err) }
	foreign, err := mailreconcile.Read(ctx, db, uuid.NewString(), "", 100)
	if err != nil || len(foreign.Rows) != 0 { t.Fatalf("tenant scope leaked: %v",err) }

	// Independent reservation requires a spare connection while the caller holds
	// its transaction. Verify bounded failure with one connection, then success with two.
	sqlDB.SetMaxOpenConns(1)
	blockedID := uuid.NewString()
	bounded, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
	err = coredatabase.RunTransaction(ctx, db, func(tx context.Context) error {
		_, err := dispatcher.Dispatch(bounded, tenant, blockedID, []byte("pool test"), send)
		return err
	})
	cancel()
	if err == nil || calls.Load() != 2 { t.Fatal("single-connection pool must fail before sending") }
	sqlDB.SetMaxOpenConns(2)
	err = coredatabase.RunTransaction(ctx, db, func(tx context.Context) error {
		_, err := dispatcher.Dispatch(tx, tenant, blockedID, []byte("pool test"), send)
		return err
	})
	if err != nil || calls.Load() != 3 { t.Fatalf("two-connection dispatch failed: %v",err) }
}
