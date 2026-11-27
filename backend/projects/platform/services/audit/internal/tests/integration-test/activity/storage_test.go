package activity_test

import (
	"context"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"grovehub/services/audit"
)

func TestAuditStorageInspection(t *testing.T) {
	dsn := os.Getenv("AUDIT_TEST_DSN")
	if dsn == "" {
		t.Skip("requires disposable PostgreSQL")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	pool, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	var database string
	if err := db.WithContext(ctx).Raw("SELECT current_database()").Scan(&database).Error; err != nil || database != "groovehub_rbac_test" {
		t.Fatal("disposable database required")
	}
	if err := audit.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := audit.InspectStorage(ctx, db, "wrong-database"); err == nil {
		t.Fatal("wrong database accepted")
	}
	report, err := audit.InspectStorage(ctx, db, database)
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"activities", "event_receipts", "rejected_events", "replay_runs"}
	if len(report) != len(expected) {
		t.Fatalf("unexpected report length: %d", len(report))
	}
	for i, table := range report {
		if table.Name != expected[i] || table.TableBytes < 0 || table.IndexBytes < 0 || table.TotalBytes < table.IndexBytes || table.EstimatedRows < 0 || table.DeadRows < 0 {
			t.Fatalf("invalid table stats: %+v", table)
		}
	}
	var count int64
	if err := db.WithContext(ctx).Raw("SELECT count(*) FROM pg_indexes WHERE schemaname = 'audit' AND indexname IN ('idx_audit_tenant_ingested','idx_audit_tenant_operation_time')").Scan(&count).Error; err != nil || count != 2 {
		t.Fatalf("missing storage indexes: %v", err)
	}
}
