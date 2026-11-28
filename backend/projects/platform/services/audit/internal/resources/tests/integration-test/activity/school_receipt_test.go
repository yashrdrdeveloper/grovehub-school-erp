package activity_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"grovehub/services/audit"
	"grovehub/services/audit/client/producer"
	operator "grovehub/services/audit/internal/app/dto/operator"
	operatorrepo "grovehub/services/audit/internal/app/repository/operator"
	operatorrepoimpl "grovehub/services/audit/internal/app/repository/operator/impl"
	operatorservice "grovehub/services/audit/internal/app/services/operator/impl"
	"os"
	"testing"
	"time"
)

func TestSchoolReceiptLedger(t *testing.T) {
	dsn := os.Getenv("AUDIT_TEST_DSN")
	if dsn == "" {
		t.Skip("requires disposable PostgreSQL")
	}
	ctx := context.Background()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	var database string
	if err = db.Raw("SELECT current_database()").Scan(&database).Error; err != nil || database != "groovehub_rbac_test" {
		t.Fatal("disposable database required")
	}
	if err = audit.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	svc := operatorservice.New(operatorrepoimpl.New(db))
	correlation := uuid.NewString()
	a := producer.SchoolRenameReceipt{ExecutionID: uuid.NewString(), SchoolID: uuid.NewString(), PreviousVersion: 1, CorrelationID: &correlation, OccurredAt: time.Now().UTC().Truncate(time.Microsecond)}
	if err = svc.StoreSchoolReceipt(ctx, a); err != nil {
		t.Fatal(err)
	}
	if err = svc.StoreSchoolReceipt(ctx, a); err != nil {
		t.Fatal("identical replay", err)
	}
	changed := a
	changed.SchoolID = uuid.NewString()
	if err = svc.StoreSchoolReceipt(ctx, changed); !errors.Is(err, operatorrepo.ErrConflict) {
		t.Fatal("different school accepted", err)
	}
	changed = a
	changed.PreviousVersion++
	if err = svc.StoreSchoolReceipt(ctx, changed); !errors.Is(err, operatorrepo.ErrConflict) {
		t.Fatal("different version accepted", err)
	}
	var n int
	if err = db.Raw("SELECT count(*) FROM audit.school_mutation_receipts WHERE execution_id=?", a.ExecutionID).Scan(&n).Error; err != nil || n != 1 {
		t.Fatal("duplicate receipt", n, err)
	}
	page, err := svc.SearchFilteredReceipts(ctx, operator.ReceiptSearchInput{SearchInput: operator.SearchInput{PageNumber: 1, RowSize: 100}})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, row := range page.Records {
		if row.ExecutionID == a.ExecutionID {
			found = row.SchoolID == a.SchoolID && row.PreviousVersion == "1" && row.CorrelationID != nil && *row.CorrelationID == correlation
		}
	}
	if !found {
		t.Fatal("missing committed receipt")
	}
	if err = db.Exec("DELETE FROM audit.school_mutation_receipts WHERE execution_id=?", a.ExecutionID).Error; err == nil {
		t.Fatal("receipt deleted")
	}
	status := producer.SchoolStatusReceipt{ExecutionID: a.ExecutionID, SchoolID: a.SchoolID, PreviousVersion: a.PreviousVersion, CorrelationID: a.CorrelationID, OccurredAt: a.OccurredAt, Operation: "DISABLE"}
	if err = svc.StoreStatusReceipt(ctx, status); !errors.Is(err, operatorrepo.ErrConflict) {
		t.Fatal("cross-operation reuse accepted", err)
	}
	status.ExecutionID = uuid.NewString()
	if err = svc.StoreStatusReceipt(ctx, status); err != nil {
		t.Fatal(err)
	}
	if err = svc.StoreStatusReceipt(ctx, status); err != nil {
		t.Fatal("status replay", err)
	}
	page, err = svc.SearchFilteredReceipts(ctx, operator.ReceiptSearchInput{SearchInput: operator.SearchInput{PageNumber: 1, RowSize: 100}})
	if err != nil {
		t.Fatal(err)
	}
	found = false
	for _, row := range page.Records {
		if row.ExecutionID == status.ExecutionID && row.Operation == "DISABLE" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing status receipt")
	}
	// Verify filtering occurs before pagination and excludes other schools.
	other := a
	other.ExecutionID = uuid.NewString()
	other.SchoolID = uuid.NewString()
	if err = svc.StoreSchoolReceipt(ctx, other); err != nil {
		t.Fatal(err)
	}
	filtered, err := svc.SearchFilteredReceipts(ctx, operator.ReceiptSearchInput{SearchInput: operator.SearchInput{PageNumber: 1, RowSize: 1}, SchoolID: a.SchoolID})
	if err != nil || len(filtered.Records) != 1 || !filtered.HasMore || filtered.Records[0].SchoolID != a.SchoolID {
		t.Fatal("wrong first filtered page", filtered, err)
	}
	filtered, err = svc.SearchFilteredReceipts(ctx, operator.ReceiptSearchInput{SearchInput: operator.SearchInput{PageNumber: 2, RowSize: 1}, SchoolID: a.SchoolID})
	if err != nil || len(filtered.Records) != 1 || filtered.HasMore || filtered.Records[0].SchoolID != a.SchoolID {
		t.Fatal("wrong second filtered page", filtered, err)
	}
	filtered, err = svc.SearchFilteredReceipts(ctx, operator.ReceiptSearchInput{SearchInput: operator.SearchInput{PageNumber: 1, RowSize: 100}, SchoolID: a.SchoolID, Operation: "DISABLE"})
	if err != nil || len(filtered.Records) != 1 || filtered.Records[0].ExecutionID != status.ExecutionID {
		t.Fatal("wrong operation filter", filtered, err)
	}
	status.Operation = "ACTIVATE"
	if err = svc.StoreStatusReceipt(ctx, status); !errors.Is(err, operatorrepo.ErrConflict) {
		t.Fatal("changed transition accepted", err)
	}
}
