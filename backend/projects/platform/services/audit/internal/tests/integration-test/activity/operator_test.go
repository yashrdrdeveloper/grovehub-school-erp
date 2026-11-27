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
 operator "grovehub/services/audit/internal/app/controller/operator/restricted"
)

func TestOperatorAuditImmutableDelivery(t *testing.T){
 dsn:=os.Getenv("AUDIT_TEST_DSN");if dsn==""{t.Skip("requires disposable PostgreSQL")}
 ctx:=context.Background()
 db,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{});if err!=nil{t.Fatal(err)}
 sqlDB,err:=db.DB();if err!=nil{t.Fatal(err)};defer sqlDB.Close()
 var database string
 if err:=db.Raw("SELECT current_database()").Scan(&database).Error;err!=nil||database!="groovehub_rbac_test"{t.Fatal("disposable database required")}
 if err=audit.Migrate(ctx,db);err!=nil{t.Fatal(err)}
 a:=producer.OperatorActivity{EventID:uuid.NewString(),ActorID:uuid.NewString(),ResourceID:uuid.NewString(),Operation:"REQUEST_APPROVED",OccurredAt:time.Now().UTC().Truncate(time.Microsecond)}
 if err=operator.Store(ctx,db,a);err!=nil{t.Fatal(err)}
 if err=operator.Store(ctx,db,a);err!=nil{t.Fatal("identical retry",err)}
 changed:=a;changed.Operation="REQUEST_REJECTED"
 if err=operator.Store(ctx,db,changed);!errors.Is(err,operator.ErrConflict){t.Fatal("conflicting occurrence accepted",err)}
 var count int
 if err=db.Raw("SELECT count(*) FROM audit.operator_activities WHERE event_id=?",a.EventID).Scan(&count).Error;err!=nil||count!=1{t.Fatal("duplicate record",count,err)}
 if err=db.Raw("SELECT count(*) FROM audit.activities WHERE source_event_id=?",a.EventID).Scan(&count).Error;err!=nil||count!=0{t.Fatal("operator record leaked into school activities",count,err)}
 page,err:=operator.Search(ctx,db,operator.SearchInput{PageNumber:1,RowSize:100})
 if err!=nil{t.Fatal(err)}
 found:=false
 for _,row:=range page.Records{if row.EventID==a.EventID{found=true}}
 if !found{t.Fatal("centralized record not found")}
 if err=db.Exec("DELETE FROM audit.operator_activities WHERE event_id=?",a.EventID).Error;err==nil{t.Fatal("audit deletion accepted")}
}
