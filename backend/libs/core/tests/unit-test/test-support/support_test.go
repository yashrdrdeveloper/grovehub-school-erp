package support_test

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"grovehub/libs/core/tests/helpers/clock"
	httphelpers "grovehub/libs/core/tests/helpers/http"
	"grovehub/libs/core/tests/mocks/controller"
	"grovehub/libs/core/tests/mocks/database"
)

func TestSharedClock(t *testing.T) {
	if got := clock.Fixed(); !got.Equal(time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC)) || got.Location() != time.UTC {
		t.Fatal(got)
	}
}

func TestHTTPFixture(t *testing.T) {
	req, w := httphelpers.Request(http.MethodPost, "/test", []byte(`{"a":1}`))
	if req.Header.Get("Content-Type") != "application/json" || req.Method != http.MethodPost || req.URL.Path != "/test" {
		t.Fatal(req)
	}
	w.WriteHeader(201)
	if w.Code != 201 {
		t.Fatal(w.Code)
	}
	req, _ = httphelpers.Request(http.MethodGet, "/test", nil)
	if req.Header.Get("Content-Type") != "" {
		t.Fatal(req.Header)
	}
}

func TestMockController(t *testing.T) {
	if ctrl := controller.New(t); ctrl == nil {
		t.Fatal("missing mock controller")
	}
}

func TestDatabaseFixture(t *testing.T) {
	db, mock, cleanup, err := database.NewPostgres()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			t.Error(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	}()
	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1))
	var value int
	if err := db.Raw("SELECT 1").Scan(&value).Error; err != nil || value != 1 {
		t.Fatalf("value=%d err=%v", value, err)
	}
	mock.ExpectClose()
}

func TestDatabaseFixtureErrors(t *testing.T) {
	sentinel := errors.New("driver unavailable")
	db, mock, cleanup, err := database.NewPostgresWith(func() (*sql.DB, sqlmock.Sqlmock, error) {
		return nil, nil, sentinel
	}, func(*sql.DB) (*gorm.DB, error) { t.Fatal("must not open after failed connection"); return nil, nil })
	if !errors.Is(err, sentinel) || db != nil || mock != nil || cleanup != nil {
		t.Fatal("connection failure not propagated")
	}
	connection, expected, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	expected.ExpectClose()
	db, mock, cleanup, err = database.NewPostgresWith(func() (*sql.DB, sqlmock.Sqlmock, error) {
		return connection, expected, nil
	}, func(*sql.DB) (*gorm.DB, error) { return nil, sentinel })
	if !errors.Is(err, sentinel) || db != nil || mock != nil || cleanup != nil {
		t.Fatal("GORM failure not propagated")
	}
	if err := expected.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
