package activity_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"grovehub/services/audit"
	"grovehub/services/audit/client/producer"
	operator "grovehub/services/audit/internal/app/dto/operator"
	operatorrepo "grovehub/services/audit/internal/app/repository/operator"
	operatorrepoimpl "grovehub/services/audit/internal/app/repository/operator/impl"
	operatorservice "grovehub/services/audit/internal/app/services/operator/impl"
)

func TestOperatorAuditImmutableDelivery(t *testing.T) {
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
	if err := db.Raw("SELECT current_database()").Scan(&database).Error; err != nil || database != "groovehub_rbac_test" {
		t.Fatal("disposable database required")
	}
	if err = audit.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	svc := operatorservice.New(operatorrepoimpl.New(db))
	correlationID := uuid.NewString()
	a := producer.OperatorActivity{CorrelationID: &correlationID, EventID: uuid.NewString(), ActorID: uuid.NewString(), ResourceID: uuid.NewString(), Operation: "REQUEST_APPROVED", OccurredAt: time.Now().UTC().Truncate(time.Microsecond)}
	if err = svc.Store(ctx, a); err != nil {
		t.Fatal(err)
	}
	if err = svc.Store(ctx, a); err != nil {
		t.Fatal("identical retry", err)
	}
	otherID := uuid.NewString()
	changedCorrelation := a
	changedCorrelation.CorrelationID = &otherID
	if err = svc.Store(ctx, changedCorrelation); !errors.Is(err, operatorrepo.ErrConflict) {
		t.Fatal("conflicting correlation accepted", err)
	}
	changed := a
	changed.Operation = "REQUEST_REJECTED"
	if err = svc.Store(ctx, changed); !errors.Is(err, operatorrepo.ErrConflict) {
		t.Fatal("conflicting occurrence accepted", err)
	}
	var count int
	if err = db.Raw("SELECT count(*) FROM audit.operator_activities WHERE event_id=?", a.EventID).Scan(&count).Error; err != nil || count != 1 {
		t.Fatal("duplicate record", count, err)
	}
	if err = db.Raw("SELECT count(*) FROM audit.activities WHERE source_event_id=?", a.EventID).Scan(&count).Error; err != nil || count != 0 {
		t.Fatal("operator record leaked into school activities", count, err)
	}
	page, err := svc.Search(ctx, operator.SearchInput{PageNumber: 1, RowSize: 100})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, row := range page.Records {
		if row.EventID == a.EventID {
			found = row.CorrelationID != nil && *row.CorrelationID == correlationID
		}
	}
	if !found {
		t.Fatal("centralized record not found")
	}
	if err = db.Exec("DELETE FROM audit.operator_activities WHERE event_id=?", a.EventID).Error; err == nil {
		t.Fatal("audit deletion accepted")
	}
}
