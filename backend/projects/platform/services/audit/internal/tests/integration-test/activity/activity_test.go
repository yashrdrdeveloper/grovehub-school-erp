package activity_test

import (
	"context"
	client "grovehub/services/audit/client/producer"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	schoolgateway "grovehub/launcher/session"
	"grovehub/services/audit"
	"grovehub/services/infrastructure"
	dto "grovehub/services/audit/internal/app/dto/activity"
	"grovehub/services/audit/internal/app/events/consumers/activity/helpers"
	consumer "grovehub/services/audit/internal/app/events/consumers/activity"
	rejection "grovehub/services/audit/internal/app/repository/rejection/impl"
	"grovehub/services/infrastructure/client/messaging"
	repo "grovehub/services/audit/internal/app/repository/activity/impl"
	"grovehub/services/audit/internal/app/repository/activity"
	service "grovehub/services/audit/internal/app/services/activity/impl"
)

type access struct { deny bool; tenantDeny bool }
func (auth access) AuthorizeProtectedRoute(ctx context.Context, pattern string) error {
	if auth.deny || (pattern!="POST /audit/protected/search" && pattern!="AUDIT tenant search") || (auth.tenantDeny && pattern=="AUDIT tenant search") { return errors.New("denied") }
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
	reused,_:=helpers.Project(ctx,eventType,raw(duplicate.SourceEventID,tenant,"other revision",99))
	if err:=repository.Insert(ctx,reused);!errors.Is(err,activity.ErrConflict) { t.Fatalf("alias reuse accepted: %v",err) }
	var aliasCount int64
	if err:=db.Raw("SELECT count(*) FROM audit.activities WHERE audit_id = ?",reused.AuditID).Scan(&aliasCount).Error;err!=nil || aliasCount!=0 { t.Fatal("conflicting alias insert did not roll back") }
	conflict,_:=helpers.Project(ctx,eventType,raw(uuid.NewString(),tenant,"changed",2))
	if err:=repository.Insert(ctx,conflict);!errors.Is(err,activity.ErrConflict) { t.Fatalf("expected collision: %v",err) }
	foreign,_:=helpers.Project(ctx,eventType,raw(uuid.NewString(),uuid.NewString(),"name",2))
	if err:=repository.Insert(ctx,foreign);err!=nil { t.Fatal(err) }
	second,_:=helpers.Project(ctx,eventType,raw(uuid.NewString(),tenant,"new",3))
	if err:=repository.Insert(ctx,second);err!=nil { t.Fatal(err) }
	if _,err:=helpers.Project(ctx,eventType,raw(uuid.NewString(),"","name",2));err==nil { t.Fatal("unscoped event accepted") }
	authCtx:=schoolgateway.WithPrincipal(ctx,schoolgateway.Principal{TenantID:tenant,AccountID:account,SessionID:"verified"})
	reader:=service.New(repository,access{tenantDeny:true})
	page,err:=reader.Search(authCtx,dto.Search{Scope:"ME",Request:coresearch.Request{PageInfo:coresearch.PageRequest{PageNumber:1,RowSize:1}}})
	if err!=nil || len(page.Items)!=1 || page.PageInfo.NextCursor=="" { t.Fatalf("first page failed: %v",err) }
	next,err:=reader.Search(authCtx,dto.Search{Scope:"ME",Request:coresearch.Request{PageInfo:coresearch.PageRequest{PageNumber:2,RowSize:1,Cursor:page.PageInfo.NextCursor}}})
	if err!=nil || len(next.Items)!=1 || next.Items[0].AuditID!=first.AuditID || next.PageInfo.NextCursor!="" { t.Fatalf("next page failed: %v",err) }
	if string(next.Items[0].Data)=="" || string(next.Items[0].Data)=="null" { t.Fatal("missing projection") }
	if page.PageInfo.TotalRecords==nil || *page.PageInfo.TotalRecords!=2 { t.Fatal("scope count incorrect") }
	if _,err:=reader.Search(authCtx,dto.Search{});err==nil { t.Fatal("self-only principal can search tenant") }
	admin:=service.New(repository,access{})
	filtered,err:=admin.Search(authCtx,dto.Search{Request:coresearch.Request{SearchText:"UPDATED",Filters:[]coresearch.Filter{{Field:"assetId",Operator:coresearch.OperatorEquals,Values:[]string{account}}},Sort:[]coresearch.Sort{{Field:"operation",Direction:coresearch.DirectionAscending},{Field:"occurredAt",Direction:coresearch.DirectionAscending}},PageInfo:coresearch.PageRequest{PageNumber:1,RowSize:1}}})
	if err!=nil || len(filtered.Items)!=1 || filtered.Items[0].AuditID!=first.AuditID || filtered.PageInfo.NextCursor=="" { t.Fatalf("filtered sort failed: %v",err) }
	if _,err:=admin.Search(authCtx,dto.Search{Request:coresearch.Request{Filters:[]coresearch.Filter{{Field:"tenantId",Operator:coresearch.OperatorEquals,Values:[]string{tenant}}}}});err==nil { t.Fatal("tenant override accepted") }
	other:=schoolgateway.WithPrincipal(ctx,schoolgateway.Principal{TenantID:tenant,AccountID:actor,SessionID:"verified"})
	empty,err:=reader.Search(other,dto.Search{Scope:"ME"})
	if err!=nil || len(empty.Items)!=0 { t.Fatal("actor-only access exposed another account") }
	var leaked bool
	if err:=db.Raw("SELECT EXISTS (SELECT 1 FROM audit.activities WHERE data::text LIKE '%never-store%')").Scan(&leaked).Error;err!=nil || leaked { t.Fatal("secret leaked into audit") }
	if _,err:=service.New(repository,access{deny:true}).Search(authCtx,dto.Search{});err==nil { t.Fatal("unauthorized search") }
	if err:=db.Exec("DELETE FROM audit.activities WHERE audit_id = ?",first.AuditID).Error;err==nil { t.Fatal("audit was mutable") }
	handler:=consumer.New(repository,rejection.New(db))
	malformed:=messaging.Message{Payload:[]byte("never-store invalid data")}
	for range 2 { if err:=handler.Handle(ctx,eventType,malformed);err!=nil { t.Fatal(err) } }
	if err:=handler.Handle(ctx,eventType,messaging.Message{Payload:raw(uuid.NewString(),tenant,"changed",2)});err!=nil { t.Fatal(err) }
	var rejected int64
	if err:=db.Raw("SELECT count(*) FROM audit.rejected_events").Scan(&rejected).Error;err!=nil || rejected!=2 { t.Fatalf("invalid rejection count %d: %v",rejected,err) }
	if err:=db.Exec("TRUNCATE audit.rejected_events").Error;err==nil { t.Fatal("rejection evidence was mutable") }
	now:=time.Now().UTC()
	occurrence:=uuid.NewString()
	failed:=client.Activity{OccurrenceID:occurrence,TenantID:tenant,Source:"identity",AssetType:"ACCOUNT",AssetID:account,Operation:"LOGIN_FAILED",OccurredAt:now}
	envelope:=func(value client.Activity) []byte {
		raw,err:=json.Marshal(map[string]any{"eventId":value.OccurrenceID,"eventType":client.ActivityOccurred,"version":1,"aggregateId":value.AssetID,"occurredAt":value.OccurredAt,"data":value})
		if err!=nil { t.Fatal(err) }
		return raw
	}
	entry,err:=helpers.Project(ctx,client.ActivityOccurred,envelope(failed))
	if err!=nil || entry.ActorID!=nil || entry.SubjectAccountID==nil || *entry.SubjectAccountID!=account { t.Fatalf("failed login projection: %v",err) }
	for range 2 { if err:=handler.Handle(ctx,client.ActivityOccurred,messaging.Message{Payload:envelope(failed)});err!=nil { t.Fatal(err) } }
	failed.ActorID=account
	if _,err:=helpers.Project(ctx,client.ActivityOccurred,envelope(failed));err==nil { t.Fatal("anonymous failure accepted authenticated actor") }
	self,err:=reader.Search(authCtx,dto.Search{Scope:"ME",Request:coresearch.Request{Filters:[]coresearch.Filter{{Field:"operation",Operator:coresearch.OperatorEquals,Values:[]string{"LOGIN_FAILED"}}}}})
	if err!=nil || len(self.Items)!=1 || self.Items[0].AuditID!=entry.AuditID { t.Fatalf("failed login self visibility: %v",err) }
	replayEvent:=uuid.NewString()
	replayPayload:=raw(replayEvent,tenant,"replay",4)
	if err:=infrastructure.NewEventStore(db).Enqueue(ctx,messaging.Event{EventID:replayEvent,Topic:"identity_account_metadata_changed",Payload:replayPayload,Metadata:map[string]string{}});err!=nil { t.Fatal(err) }
	source:=infrastructure.NewReplaySource(db)
	request:=audit.ReplayRequest{TenantID:tenant,EventID:replayEvent}
	if _,err:=audit.Replay(ctx,db,source,request);err==nil { t.Fatal("unpublished replay accepted") }
	if err:=db.Exec("UPDATE infrastructure.messaging SET status = 'SENT',published_at=now() WHERE event_id = ?",replayEvent).Error;err!=nil { t.Fatal(err) }
	preview,err:=audit.Replay(ctx,db,source,request)
	if err!=nil || preview.Applied || preview.AuditID=="" { t.Fatalf("replay preview: %v",err) }
	var count int64
	if err:=db.Raw("SELECT count(*) FROM audit.activities WHERE audit_id = ?",preview.AuditID).Scan(&count).Error;err!=nil || count!=0 { t.Fatal("preview mutated data") }
	request.Apply=true
	request.OperatorLabel="integration-test"
	wrong:=request
	wrong.TenantID=uuid.NewString()
	if _,err:=audit.Replay(ctx,db,source,wrong);err==nil { t.Fatal("cross-tenant replay accepted") }
	for range 2 {
		applied,err:=audit.Replay(ctx,db,source,request)
		if err!=nil || !applied.Applied || applied.AuditID!=preview.AuditID || applied.ReplayID=="" { t.Fatalf("replay apply: %v",err) }
	}
	if err:=db.Raw("SELECT count(*) FROM audit.activities WHERE audit_id = ?",preview.AuditID).Scan(&count).Error;err!=nil || count!=1 { t.Fatal("replay duplicated activity") }
	if err:=db.Raw("SELECT count(*) FROM audit.replay_runs WHERE source_event_id = ?",replayEvent).Scan(&count).Error;err!=nil || count!=2 { t.Fatal("missing replay receipts") }
	if err:=db.Exec("DELETE FROM audit.replay_runs WHERE source_event_id = ?",replayEvent).Error;err==nil { t.Fatal("mutable replay evidence") }
}
