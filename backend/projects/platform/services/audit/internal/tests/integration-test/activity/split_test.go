package activity_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/launcher/bootstrap"
	"grovehub/services/audit"
	client "grovehub/services/audit/client/producer"
	"grovehub/services/infrastructure"
)

func TestSeparateAuditDatabase(t *testing.T) {
	if os.Getenv("AUDIT_SPLIT_TEST_DSN")=="" || os.Getenv("AUDIT_TEST_DSN")=="" { t.Skip("requires two disposable PostgreSQL databases") }
	t.Setenv("DATABASE_URL",os.Getenv("AUDIT_TEST_DSN"))
	t.Setenv("AUDIT_DATABASE_URL","")
	ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second)
	defer cancel()
	business,err:=bootstrap.NewPostgresDB(ctx)
	if err!=nil { t.Fatal(err) }
	sqlDB,err:=business.DB()
	if err!=nil { t.Fatal(err) }
	defer sqlDB.Close()
	checkDatabase(t,business,"groovehub_rbac_test")
	shared,release,err:=bootstrap.NewAuditPostgresDB(ctx,business)
	if err!=nil || shared!=business { t.Fatalf("shared pool: %v",err) }
	release()
	if err:=sqlDB.PingContext(ctx);err!=nil { t.Fatal("borrowed pool closed") }
	t.Setenv("AUDIT_DATABASE_URL",os.Getenv("AUDIT_SPLIT_TEST_DSN"))
	target,closeTarget,err:=bootstrap.NewAuditPostgresDB(ctx,business)
	if err!=nil { t.Fatal(err) }
	defer closeTarget()
	checkDatabase(t,target,"groovehub_audit_test")
	if target==business { t.Fatal("Audit did not open independent pool") }
	if err:=audit.Migrate(ctx,target);err!=nil { t.Fatal(err) }
	var hasOutbox bool
	if err:=target.Raw("SELECT to_regclass('infrastructure.messaging') IS NOT NULL").Scan(&hasOutbox).Error;err!=nil || hasOutbox { t.Fatal("Audit database must not contain business outbox") }
	tenant,account,event:=uuid.NewString(),uuid.NewString(),uuid.NewString()
	producer:=audit.NewProducer(infrastructure.NewEventStore(business))
	activity:=client.Activity{OccurrenceID:event,TenantID:tenant,Source:"identity",AssetType:"ACCOUNT",AssetID:account,ActorID:account,Operation:"LOGIN_SUCCEEDED",OccurredAt:time.Now().UTC()}
	rollback:=errors.New("injected rollback")
	if err:=coredatabase.RunTransaction(ctx,business,func(tx context.Context) error {
		if err:=producer.Record(tx,activity);err!=nil { return err }
		return rollback
	});!errors.Is(err,rollback) { t.Fatalf("producer transaction rollback: %v",err) }
	var count int64
	if err:=business.Raw("SELECT count(*) FROM infrastructure.messaging WHERE event_id = ?",event).Scan(&count).Error;err!=nil || count!=0 { t.Fatal("rolled-back outbox persisted") }
	if err:=coredatabase.RunTransaction(ctx,business,func(tx context.Context) error { return producer.Record(tx,activity) });err!=nil { t.Fatal(err) }
	if err:=business.Exec("UPDATE infrastructure.messaging SET status='SENT',published_at=now() WHERE event_id = ?",event).Error;err!=nil { t.Fatal(err) }
	result,err:=audit.Replay(ctx,target,infrastructure.NewReplaySource(business),audit.ReplayRequest{TenantID:tenant,EventID:event,Apply:true,OperatorLabel:"split-test"})
	if err!=nil || !result.Applied { t.Fatalf("split replay: %v",err) }
	if err:=target.Raw("SELECT count(*) FROM audit.activities WHERE source_event_id = ?",event).Scan(&count).Error;err!=nil || count!=1 { t.Fatal("missing target activity") }
	if err:=business.Raw("SELECT count(*) FROM audit.activities WHERE source_event_id = ?",event).Scan(&count).Error;err!=nil || count!=0 { t.Fatal("activity written to business database") }
}

func checkDatabase(t *testing.T,db *gorm.DB,expected string) {
	t.Helper()
	var name string
	if err:=db.Raw("SELECT current_database()").Scan(&name).Error;err!=nil || name!=expected { t.Fatal("disposable database required") }
}
