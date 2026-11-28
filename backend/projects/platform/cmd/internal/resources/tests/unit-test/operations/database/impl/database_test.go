package impl_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"grovehub/cmd/internal/src/operations/database/impl"
)

func TestDatabase(t *testing.T) {
	ctx := context.Background()
	failure := errors.New("open failed")
	db := impl.NewWithOpeners(func(context.Context) (*gorm.DB, error) {
		return nil, failure
	}, nil)
	if got, closeDB, err := db.Business(ctx); got != nil || closeDB != nil || !errors.Is(err, failure) {
		t.Fatal("business opener error not propagated")
	}

	db = impl.NewWithOpeners(func(context.Context) (*gorm.DB, error) {
		return &gorm.DB{Config: &gorm.Config{}}, nil
	}, nil)
	if got, closeDB, err := db.Business(ctx); got != nil || closeDB != nil || err == nil {
		t.Fatal("missing underlying pool not rejected")
	}

	raw, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	pool, err := gorm.Open(postgres.New(postgres.Config{Conn: raw}), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}

	auditClosed := false
	db = impl.NewWithOpeners(func(context.Context) (*gorm.DB, error) {
		return pool, nil
	}, func(_ context.Context, shared *gorm.DB) (*gorm.DB, func(), error) {
		if shared != pool {
			t.Fatal("shared database not forwarded")
		}

		return pool, func() { auditClosed = true }, nil
	})

	business, closeBusiness, err := db.Business(ctx)
	if err != nil || business != pool {
		t.Fatal("business pool not returned", err)
	}
	audit, closeAudit, err := db.Audit(ctx, business)
	if err != nil || audit != pool {
		t.Fatal("audit pool not returned", err)
	}

	mock.ExpectQuery("SELECT current_database").WillReturnRows(sqlmock.NewRows([]string{"current_database"}).AddRow("fixture"))
	if name, err := db.CurrentName(ctx, pool); name != "fixture" || err != nil {
		t.Fatalf("database name=%q error=%v", name, err)
	}
	mock.ExpectQuery("SELECT current_database").WillReturnError(failure)
	if _, err := db.CurrentName(ctx, pool); !errors.Is(err, failure) {
		t.Fatal("name lookup failure not returned", err)
	}

	closeAudit()
	mock.ExpectClose()
	closeBusiness()
	if !auditClosed {
		t.Fatal("audit cleanup not called")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
