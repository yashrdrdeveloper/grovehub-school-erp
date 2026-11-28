package support_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/tests/helpers/clock"
	httphelpers "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/tests/helpers/http"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/tests/mocks/controller"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/tests/mocks/database"
)

func TestSharedClock(t *testing.T) {
	if got:=clock.Fixed(); !got.Equal(time.Date(2026,time.November,14,0,0,0,0,time.UTC)) || got.Location()!=time.UTC { t.Fatal(got) }
}

func TestHTTPFixture(t *testing.T) {
	req,w:=httphelpers.Request(http.MethodPost,"/test",[]byte(`{"a":1}`))
	if req.Header.Get("Content-Type")!="application/json" || req.Method!=http.MethodPost || req.URL.Path!="/test" { t.Fatal(req) }
	w.WriteHeader(201)
	if w.Code!=201 { t.Fatal(w.Code) }
	req,_=httphelpers.Request(http.MethodGet,"/test",nil)
	if req.Header.Get("Content-Type")!="" { t.Fatal(req.Header) }
}

func TestMockController(t *testing.T) {
	if ctrl:=controller.New(t);ctrl==nil { t.Fatal("missing mock controller") }
}

func TestDatabaseFixture(t *testing.T) {
	db,mock,cleanup,err:=database.NewPostgres()
	if err!=nil { t.Fatal(err) }
	defer func(){if err:=cleanup();err!=nil{t.Error(err)};if err:=mock.ExpectationsWereMet();err!=nil{t.Error(err)}}()
	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1))
	var value int
	if err:=db.Raw("SELECT 1").Scan(&value).Error;err!=nil || value!=1 { t.Fatalf("value=%d err=%v",value,err) }
	mock.ExpectClose()
}
