package delivery

import (
 "context"
 "crypto/ed25519"
 "database/sql"
 "encoding/base64"
 "errors"
 "log/slog"
 "net/http"
 "net/url"
 "strings"
 "time"

 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/generated/audit"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/services"
)
const statusPath="/audit/restricted/school-status-receipts"
const statusScope="internal.audit.school.status-receipt.write"
const statusIssuer="school-registry-status-audit"
type StatusPublisher func(context.Context,sdk.SchoolStatusReceipt)(int,error)
type StatusWorker struct{db *sql.DB;schoolID string;publish StatusPublisher}
func NewStatus(db *sql.DB,school string,publish StatusPublisher)*StatusWorker{return &StatusWorker{db:db,schoolID:school,publish:publish}}
func StatusClient(c Config)(StatusPublisher,error){
 if c.Operation!="status"{return nil,errors.New("invalid status publisher operation")}
 origin,err:=url.Parse(c.Origin)
 if err!=nil||origin.Hostname()==""||origin.User!=nil||origin.RawQuery!=""||origin.ForceQuery||origin.Fragment!=""||origin.RawPath!=""||(origin.Path!=""&&origin.Path!="/")||(origin.Scheme!="https"&&!(origin.Scheme=="http"&&(origin.Hostname()=="127.0.0.1"||origin.Hostname()=="::1"))){return nil,errors.New("invalid receipt origin")}
 key,err:=base64.RawURLEncoding.DecodeString(c.PrivateKey)
 if err!=nil||len(key)!=ed25519.PrivateKeySize||!services.ValidID(c.SchoolID)||c.KeyID==""{return nil,errors.New("invalid receipt credential")}
 base:=http.DefaultTransport.(*http.Transport).Clone();base.Proxy=nil
 t:=&core.WorkloadTransport{Origin:origin,KeyID:c.KeyID,Issuer:statusIssuer,Subject:statusIssuer+":"+c.SchoolID,Audience:"platform",PrivateKey:ed25519.PrivateKey(key),Base:base,Policy:func(r *http.Request)(core.WorkloadClaims,error){
  if r==nil||r.URL==nil||r.Method!="POST"||r.URL.Path!=statusPath||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return core.WorkloadClaims{},core.ErrWorkloadDenied}
  return core.WorkloadClaims{Scopes:[]string{statusScope}},nil
 }}
 client,err:=core.NewWorkloadHTTPClient(t);if err!=nil{return nil,err}
 client.Transport=bounded{client.Transport}
 cfg:=sdk.NewConfiguration();cfg.Servers=sdk.ServerConfigurations{{URL:strings.TrimRight(origin.String(),"/")}};cfg.HTTPClient=client
 api:=sdk.NewAPIClient(cfg)
 return func(ctx context.Context,a sdk.SchoolStatusReceipt)(int,error){
  if a.SchoolId!=c.SchoolID{return 0,core.ErrWorkloadDenied}
  resp,err:=api.DefaultAPI.RecordSchoolStatusReceipt(ctx).SchoolStatusReceipt(a).Execute()
  if resp==nil{return 0,errors.New("receipt unavailable")}
  if resp.Body!=nil{defer resp.Body.Close()}
  if resp.StatusCode==204&&err!=nil{return 0,errors.New("invalid receipt acknowledgment")}
  return resp.StatusCode,nil
 },nil
}
func(w *StatusWorker) Step(ctx context.Context)(bool,error){
 ctx,cancel:=context.WithTimeout(ctx,10*time.Second);defer cancel()
 tx,err:=w.db.BeginTx(ctx,nil);if err!=nil{return false,err};defer tx.Rollback()
 var a sdk.SchoolStatusReceipt
 err=tx.QueryRowContext(ctx,"SELECT e.execution_id,e.school_id,e.previous_version,e.correlation_id,e.occurred_at,CASE WHEN e.new_status='ACTIVE' THEN 'ACTIVATE' ELSE 'DISABLE' END FROM school_registry.status_outbox o JOIN school_registry.status_events e ON e.execution_id=o.execution_id WHERE e.school_id=$1 AND o.delivered_at IS NULL AND o.blocked_at IS NULL AND o.next_attempt_at<=now() ORDER BY o.next_attempt_at,o.execution_id LIMIT 1 FOR UPDATE OF o SKIP LOCKED",w.schoolID).Scan(&a.ExecutionId,&a.SchoolId,&a.PreviousVersion,&a.CorrelationId,&a.OccurredAt,&a.Operation)
 if errors.Is(err,sql.ErrNoRows){return false,nil};if err!=nil{return false,err}
 probe,stop:=context.WithTimeout(ctx,5*time.Second)
 status,sendErr:=w.publish(probe,a);stop()
 switch{
 case sendErr==nil&&status==204:
  _,err=tx.ExecContext(ctx,"UPDATE school_registry.status_outbox SET delivered_at=now(),attempts=attempts+1 WHERE execution_id=$1",a.ExecutionId)
 case sendErr==nil&&(status==400||status==409):
  _,err=tx.ExecContext(ctx,"UPDATE school_registry.status_outbox SET blocked_at=now(),attempts=attempts+1 WHERE execution_id=$1",a.ExecutionId)
  slog.Warn("school receipt blocked","execution_id",a.ExecutionId,"status",status)
 default:
  _,err=tx.ExecContext(ctx,"UPDATE school_registry.status_outbox SET next_attempt_at=now()+interval '30 seconds',attempts=attempts+1 WHERE execution_id=$1",a.ExecutionId)
 }
 if err!=nil{return false,err};return true,tx.Commit()
}
func(w *StatusWorker)Run(ctx context.Context){
 for ctx.Err()==nil{
  worked,err:=w.Step(ctx);if err!=nil&&ctx.Err()==nil{slog.Warn("school receipt relay failed")}
  if worked&&err==nil{continue}
  select{case <-ctx.Done():return;case <-time.After(time.Second):}
 }
}
