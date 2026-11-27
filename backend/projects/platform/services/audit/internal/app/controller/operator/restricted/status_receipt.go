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
const StatusReceiptPattern="POST /audit/restricted/school-status-receipts"
const StatusReceiptScope="internal.audit.school.status-receipt.write"
const StatusReceiptIssuer="school-registry-status-audit"
func StatusReceiptHandler(db *gorm.DB)http.HandlerFunc{
 return func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  claims,ok:=core.WorkloadFrom(r.Context())
  if !ok||len(claims.Scopes)!=1||claims.Scopes[0]!=StatusReceiptScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{w.WriteHeader(403);return}
  var a producer.SchoolStatusReceipt;var extra any
  d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096));d.DisallowUnknownFields()
  if d.Decode(&a)!=nil||d.Decode(&extra)!=io.EOF||!validID(a.ExecutionID)||!validID(a.SchoolID)||(a.Operation!="ACTIVATE"&&a.Operation!="DISABLE")||a.PreviousVersion<1||(a.CorrelationID!=nil&&!validID(*a.CorrelationID))||a.OccurredAt.IsZero()||a.OccurredAt.After(time.Now().Add(5*time.Minute)){w.WriteHeader(400);return}
  if claims.Subject!=StatusReceiptIssuer+":"+a.SchoolID{w.WriteHeader(403);return}
  a.OccurredAt=a.OccurredAt.UTC().Truncate(time.Microsecond)
  ctx,cancel:=context.WithTimeout(r.Context(),5*time.Second);defer cancel()
  err:=StoreStatusReceipt(ctx,db,a)
  if errors.Is(err,ErrConflict){w.WriteHeader(409);return}
  if err!=nil{w.WriteHeader(503);return};w.WriteHeader(204)
 }
}
