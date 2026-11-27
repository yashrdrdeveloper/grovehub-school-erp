package activity_test

import(
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
func TestSchoolReceiptLedger(t *testing.T){
 dsn:=os.Getenv("AUDIT_TEST_DSN");if dsn==""{t.Skip("requires disposable PostgreSQL")}
 ctx:=context.Background()
 db,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{});if err!=nil{t.Fatal(err)}
 sqlDB,err:=db.DB();if err!=nil{t.Fatal(err)};defer sqlDB.Close()
 var database string
 if err=db.Raw("SELECT current_database()").Scan(&database).Error;err!=nil||database!="groovehub_rbac_test"{t.Fatal("disposable database required")}
 if err=audit.Migrate(ctx,db);err!=nil{t.Fatal(err)}
 correlation:=uuid.NewString()
 a:=producer.SchoolRenameReceipt{ExecutionID:uuid.NewString(),SchoolID:uuid.NewString(),PreviousVersion:1,CorrelationID:&correlation,OccurredAt:time.Now().UTC().Truncate(time.Microsecond)}
 if err=operator.StoreSchoolReceipt(ctx,db,a);err!=nil{t.Fatal(err)}
 if err=operator.StoreSchoolReceipt(ctx,db,a);err!=nil{t.Fatal("identical replay",err)}
 changed:=a;changed.SchoolID=uuid.NewString()
 if err=operator.StoreSchoolReceipt(ctx,db,changed);!errors.Is(err,operator.ErrConflict){t.Fatal("different school accepted",err)}
 changed=a;changed.PreviousVersion++
 if err=operator.StoreSchoolReceipt(ctx,db,changed);!errors.Is(err,operator.ErrConflict){t.Fatal("different version accepted",err)}
 var n int
 if err=db.Raw("SELECT count(*) FROM audit.school_mutation_receipts WHERE execution_id=?",a.ExecutionID).Scan(&n).Error;err!=nil||n!=1{t.Fatal("duplicate receipt",n,err)}
 page,err:=operator.SearchReceipts(ctx,db,operator.SearchInput{PageNumber:1,RowSize:100})
 if err!=nil{t.Fatal(err)}
 found:=false
 for _,row:=range page.Records{if row.ExecutionID==a.ExecutionID{found=row.SchoolID==a.SchoolID&&row.PreviousVersion=="1"&&row.CorrelationID!=nil&&*row.CorrelationID==correlation}}
 if !found{t.Fatal("missing committed receipt")}
 if err=db.Exec("DELETE FROM audit.school_mutation_receipts WHERE execution_id=?",a.ExecutionID).Error;err==nil{t.Fatal("receipt deleted")}
}
