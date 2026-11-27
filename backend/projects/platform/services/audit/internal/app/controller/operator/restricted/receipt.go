package restricted

import (
 "context"
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "time"
 "gorm.io/gorm"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 "grovehub/services/audit/client/producer"
)
const ReceiptPattern="POST /audit/restricted/school-rename-receipts"
const ReceiptScope="internal.audit.school.rename-receipt.write"
const ReceiptIssuer="school-registry-audit"
func StoreSchoolReceipt(ctx context.Context,db *gorm.DB,a producer.SchoolRenameReceipt)error{
 return storeReceipt(ctx,db,a,"RENAME")
}
func StoreStatusReceipt(ctx context.Context,db *gorm.DB,a producer.SchoolStatusReceipt)error{
 if a.Operation!="ACTIVATE"&&a.Operation!="DISABLE"{return errors.New("invalid lifecycle operation")}
 return storeReceipt(ctx,db,producer.SchoolRenameReceipt{ExecutionID:a.ExecutionID,SchoolID:a.SchoolID,PreviousVersion:a.PreviousVersion,CorrelationID:a.CorrelationID,OccurredAt:a.OccurredAt},a.Operation)
}
func storeReceipt(ctx context.Context,db *gorm.DB,a producer.SchoolRenameReceipt,operation string)error{
 return db.WithContext(ctx).Transaction(func(tx *gorm.DB)error{
  result:=tx.Exec("INSERT INTO audit.school_mutation_receipts(execution_id,school_id,previous_version,correlation_id,occurred_at,operation) VALUES(?,?,?,?,?,?) ON CONFLICT(execution_id) DO NOTHING",a.ExecutionID,a.SchoolID,a.PreviousVersion,a.CorrelationID,a.OccurredAt,operation)
  if result.Error!=nil{return result.Error};if result.RowsAffected==1{return nil}
  var same bool
  if err:=tx.Raw("SELECT EXISTS(SELECT 1 FROM audit.school_mutation_receipts WHERE execution_id=? AND school_id=? AND previous_version=? AND correlation_id IS NOT DISTINCT FROM CAST(? AS uuid) AND occurred_at=? AND operation=?)",a.ExecutionID,a.SchoolID,a.PreviousVersion,a.CorrelationID,a.OccurredAt,operation).Scan(&same).Error;err!=nil{return err}
  if !same{return ErrConflict};return nil
 })
}
func ReceiptHandler(db *gorm.DB)http.HandlerFunc{
 return func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  claims,ok:=core.WorkloadFrom(r.Context())
  if !ok||len(claims.Scopes)!=1||claims.Scopes[0]!=ReceiptScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{w.WriteHeader(403);return}
  var a producer.SchoolRenameReceipt;var extra any
  d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096));d.DisallowUnknownFields()
  if d.Decode(&a)!=nil||d.Decode(&extra)!=io.EOF||!validID(a.ExecutionID)||!validID(a.SchoolID)||a.PreviousVersion<1||(a.CorrelationID!=nil&&!validID(*a.CorrelationID))||a.OccurredAt.IsZero()||a.OccurredAt.After(time.Now().Add(5*time.Minute)){w.WriteHeader(400);return}
  if claims.Subject!=ReceiptIssuer+":"+a.SchoolID{w.WriteHeader(403);return}
  a.OccurredAt=a.OccurredAt.UTC().Truncate(time.Microsecond)
  ctx,cancel:=context.WithTimeout(r.Context(),5*time.Second);defer cancel()
  err:=StoreSchoolReceipt(ctx,db,a)
  if errors.Is(err,ErrConflict){w.WriteHeader(409);return}
  if err!=nil{w.WriteHeader(503);return};w.WriteHeader(204)
 }
}
