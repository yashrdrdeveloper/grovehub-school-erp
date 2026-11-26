package activity_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"grovehub/libs/schoolgateway"
	"grovehub/services/audit"
	dto "grovehub/services/audit/internal/app/dto/activity"
	"grovehub/services/audit/internal/app/events/consumers/activity/helpers"
	repo "grovehub/services/audit/internal/app/repository/activity/impl"
	"grovehub/services/audit/internal/app/repository/activity"
	service "grovehub/services/audit/internal/app/services/activity/impl"
)

type access struct { deny bool }
func (auth access) AuthorizeProtectedRoute(ctx context.Context, pattern string) error {
	if auth.deny || pattern!="POST /audit/protected/search" { return errors.New("denied") }
	return nil
}

func TestAuditProjectionAndPersistence(t *testing.T) {
	dsn:=os.Getenv("AUDIT_TEST_DSN")
	if dsn=="" { t.Skip("requires disposable PostgreSQL") }
	ctx:=context.Background()
	db,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil { t.Fatal(err) }
	sqlDB,err:=db.DB()
	if err!=nil { t.Fatal(err) }
	defer sqlDB.Close()
	var database string
	if err:=db.Raw("SELECT current_database()").Scan(&database).Error;err!=nil || database!="groovehub_rbac_test" { t.Fatal("disposable database required") }
	if err:=audit.Migrate(ctx,db);err!=nil { t.Fatal(err) }
	tenant,account,actor:=uuid.NewString(),uuid.NewString(),uuid.NewString()
	eventType:="identity.account.metadata.changed.v1"
	raw:=func(eventID string, tenantID string, name string, version int) []byte {
		value,_:=json.Marshal(map[string]any{"eventId":eventID,"eventType":eventType,"version":1,"aggregateId":account,"occurredAt":time.Now().UTC(),"data":map[string]any{"tenantId":tenantID,"actorId":actor,"operation":"UPDATED","password":"never-store","new":map[string]any{"accountId":account,"username":name,"version":version,"token":"never-store"}}})
		return value
	}
	repository:=repo.New(db)
	first,err:=helpers.Project(ctx,eventType,raw(uuid.NewString(),tenant,"name",2))
	if err!=nil { t.Fatal(err) }
	duplicate,err:=helpers.Project(ctx,eventType,raw(uuid.NewString(),tenant,"name",2))
	if err!=nil || first.AuditID!=duplicate.AuditID { t.Fatal("logical replay generated different ID") }
	if err:=repository.Insert(ctx,first);err!=nil { t.Fatal(err) }
	if err:=repository.Insert(ctx,duplicate);err!=nil { t.Fatal(err) }
	conflict,_:=helpers.Project(ctx,eventType,raw(uuid.NewString(),tenant,"changed",2))
	if err:=repository.Insert(ctx,conflict);!errors.Is(err,activity.ErrConflict) { t.Fatalf("expected collision: %v",err) }
	foreign,_:=helpers.Project(ctx,eventType,raw(uuid.NewString(),uuid.NewString(),"name",2))
	if err:=repository.Insert(ctx,foreign);err!=nil { t.Fatal(err) }
	second,_:=helpers.Project(ctx,eventType,raw(uuid.NewString(),tenant,"new",3))
	if err:=repository.Insert(ctx,second);err!=nil { t.Fatal(err) }
	if _,err:=helpers.Project(ctx,eventType,raw(uuid.NewString(),"","name",2));err==nil { t.Fatal("unscoped event accepted") }
	authCtx:=schoolgateway.WithPrincipal(ctx,schoolgateway.Principal{TenantID:tenant,AccountID:actor,SessionID:"verified"})
	reader:=service.New(repository,access{})
	page,err:=reader.Search(authCtx,dto.Search{Scope:"ME",Limit:1})
	if err!=nil || len(page.Items)!=1 || page.NextCursor==nil { t.Fatalf("first page failed: %v",err) }
	next,err:=reader.Search(authCtx,dto.Search{Scope:"ME",Limit:1,Cursor:page.NextCursor})
	if err!=nil || len(next.Items)!=1 || next.Items[0].AuditID!=first.AuditID || next.NextCursor!=nil { t.Fatalf("next page failed: %v",err) }
	if string(next.Items[0].Data)=="" || string(next.Items[0].Data)=="null" { t.Fatal("missing projection") }
	var leaked bool
	if err:=db.Raw("SELECT EXISTS (SELECT 1 FROM audit.activities WHERE data::text LIKE '%never-store%')").Scan(&leaked).Error;err!=nil || leaked { t.Fatal("secret leaked into audit") }
	if _,err:=service.New(repository,access{deny:true}).Search(authCtx,dto.Search{});err==nil { t.Fatal("unauthorized search") }
	if err:=db.Exec("DELETE FROM audit.activities WHERE audit_id = ?",first.AuditID).Error;err==nil { t.Fatal("audit was mutable") }
}
