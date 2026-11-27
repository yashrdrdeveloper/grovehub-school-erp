package restricted

import (
 "context"
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "time"

 "github.com/google/uuid"
 "gorm.io/gorm"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 "grovehub/services/audit/client/producer"
)

const Scope="internal.audit.operator.write"
const Pattern="POST /audit/restricted/operator-activities"

func validID(value string)bool{
 id,err:=uuid.Parse(value)
 return err==nil&&id!=uuid.Nil&&id.String()==value
}
func validOperation(op string)bool{
 switch op{
 case "ENVIRONMENT_CREATED","TEMPLATE_CREATED","REQUEST_SUBMITTED","REQUEST_AUTO_APPROVED","REQUEST_APPROVED","REQUEST_REJECTED","REQUEST_CANCELLED","EXECUTION_STARTED","EXECUTION_SUCCEEDED","EXECUTION_FAILED","EXECUTION_UNKNOWN","EXECUTION_ABANDONED","EXECUTION_RECONCILED","OPERATOR_CREATED","OPERATOR_ACTIVE","OPERATOR_DISABLED","AUDIT_DELIVERY_REQUEUED","OPERATOR_BOOTSTRAPPED","AUTH_LOGIN_SUCCEEDED","AUTH_LOGOUT_SUCCEEDED":
  return true
 }
 return false
}
func valid(a producer.OperatorActivity)bool{
 return validID(a.EventID)&&validID(a.ActorID)&&validID(a.ResourceID)&&validOperation(a.Operation)&&(a.CorrelationID==nil||validID(*a.CorrelationID))&&!a.OccurredAt.IsZero()&&!a.OccurredAt.After(time.Now().Add(5*time.Minute))
}
var ErrConflict=errors.New("operator occurrence conflict")

// Store acknowledges only after durable insertion. Duplicate retries must have
// exactly the same immutable identity and semantic fields.
func Store(ctx context.Context,db *gorm.DB,a producer.OperatorActivity)error{
 return db.WithContext(ctx).Transaction(func(tx *gorm.DB)error{
  result:=tx.Exec("INSERT INTO audit.operator_activities(event_id,source,actor_id,resource_id,operation,occurred_at,correlation_id) VALUES(?,'internal-tool',?,?,?,?,?) ON CONFLICT(event_id) DO NOTHING",a.EventID,a.ActorID,a.ResourceID,a.Operation,a.OccurredAt.UTC(),a.CorrelationID)
  if result.Error!=nil{return result.Error}
  if result.RowsAffected==1{return nil}
  var identical bool
  if err:=tx.Raw("SELECT EXISTS(SELECT 1 FROM audit.operator_activities WHERE event_id=? AND source='internal-tool' AND actor_id=? AND resource_id=? AND operation=? AND occurred_at=? AND correlation_id IS NOT DISTINCT FROM CAST(? AS uuid))",a.EventID,a.ActorID,a.ResourceID,a.Operation,a.OccurredAt.UTC(),a.CorrelationID).Scan(&identical).Error;err!=nil{return err}
  if !identical{return ErrConflict};return nil
 })
}
func Handler(db *gorm.DB)http.HandlerFunc{
 return func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  claims,ok:=core.WorkloadFrom(r.Context())
  if !ok||claims.Subject!="internal-tool-audit"||len(claims.Scopes)!=1||claims.Scopes[0]!=Scope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{w.WriteHeader(403);return}
  var a producer.OperatorActivity
  d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096));d.DisallowUnknownFields()
  var extra any
  if d.Decode(&a)!=nil||d.Decode(&extra)!=io.EOF||!valid(a){w.WriteHeader(400);return}
  // PostgreSQL timestamps have microsecond precision; normalize for retries.
  a.OccurredAt=a.OccurredAt.UTC().Truncate(time.Microsecond)
  ctx,cancel:=context.WithTimeout(r.Context(),5*time.Second);defer cancel()
  err:=Store(ctx,db,a)
  if errors.Is(err,ErrConflict){w.WriteHeader(409);return}
  if err!=nil{w.WriteHeader(503);return}
  w.WriteHeader(204)
 }
}
