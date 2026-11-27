package workflow

import (
 "context"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/correlation"
 "database/sql"
 "encoding/json"
 "errors"
 "regexp"
 "strings"
 "time"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

var ErrExecutionUnavailable=errors.New("execution binding unavailable")
var segmentPattern=regexp.MustCompile("^[a-zA-Z0-9_.-]{1,128}$")
type ExecutionInput struct{Version int64 `json:"version"`}
type BindingKey struct{EnvironmentCode,TemplateID string}
type Invocation struct{
 ExecutionID string
 RequestID string
 SchoolID *string
 Method,Path string
 Headers map[string]string
 Input json.RawMessage
}
// Adapter is compiled/configured by the deployment owner, not registered by
// an operator. It must independently enforce destination, operation and tenant.
// No generic URL transport is provided by the workflow engine.
type Adapter func(context.Context,Invocation)(int,error)

func NewWithBindings(db *sql.DB,bindings map[BindingKey]Adapter)*Service{
 s:=New(db);s.bindings=make(map[BindingKey]Adapter,len(bindings))
 for k,v:=range bindings{if v!=nil{s.bindings[k]=v}}
 return s
}
func resolve(t entity.Template,a entity.Approval)(Invocation,error){
 if !validInput(a.Input,t.BodySchema){return Invocation{},ErrInvalid}
 var values map[string]json.RawMessage
 if err:=json.Unmarshal(a.Input,&values);err!=nil{return Invocation{},ErrInvalid}
 lookup:=func(m string)(string,error){
  parts:=placeholderPattern.FindStringSubmatch(m);if len(parts)!=2{return "",ErrInvalid}
  var value string
  if json.Unmarshal(values[parts[1]],&value)!=nil||!segmentPattern.MatchString(value)||value=="."||value==".."{return "",ErrInvalid}
  return value,nil
 }
 path:=t.Endpoint
 for _,m:=range placeholderPattern.FindAllString(path,-1){value,err:=lookup(m);if err!=nil{return Invocation{},err};path=strings.ReplaceAll(path,m,value)}
 headers:=map[string]string{}
 for key,value:=range t.Headers{
  if key=="Idempotency-Key"{v,err:=lookup(value);if err!=nil{return Invocation{},err};value=v}
  headers[key]=value
 }
 return Invocation{RequestID:a.ID,SchoolID:a.SchoolID,Method:t.Method,Path:path,Headers:headers,Input:append(json.RawMessage(nil),a.Input...)},nil
}
func(s *Service) Execute(ctx context.Context,actor entity.Operator,id string,in ExecutionInput)(entity.Execution,error){
 ctx=correlation.Ensure(ctx)
 if !validID(id)||in.Version<1{return entity.Execution{},ErrInvalid}
 a,err:=s.repo.Get(ctx,actor,id);if err!=nil{return entity.Execution{},err}
 if a.Status!="APPROVED"||a.Version!=in.Version{return entity.Execution{},ErrConflict}
 t,err:=s.repo.Template(ctx,a.TemplateID);if err!=nil{return entity.Execution{},err}
 env,err:=s.repo.Environment(ctx,a.EnvironmentID);if err!=nil{return entity.Execution{},err}
 call,err:=resolve(t,a);if err!=nil{return entity.Execution{},err}
 adapter,err:=s.executionAdapter(env.Code,t,a,actor,call);if err!=nil{return entity.Execution{},err}
 if a.SchoolID!=nil{if err:=s.checkDirectoryForBinding(ctx,*a.SchoolID,s.schoolBindings[SchoolBindingKey{env.Code,t.ID,*a.SchoolID}]);err!=nil{return entity.Execution{},err}}
 e,err:=s.repo.ClaimExecution(ctx,actor,id,uuid.NewString(),in.Version);if err!=nil{return e,err}
 call.ExecutionID=e.ID
 // At-most-one dispatch claim: ambiguous outcomes are never retried automatically.
 dispatch,cancel:=context.WithTimeout(context.WithoutCancel(ctx),15*time.Second)
 status,callErr:=adapter(dispatch,call);cancel()
 outcome:="UNKNOWN";var response *int
 if callErr==nil&&status>=200&&status<=599{
  response=&status
  if status<300{outcome="SUCCEEDED"}else if status<500{outcome="FAILED"}
 }
 finish,stop:=context.WithTimeout(context.WithoutCancel(ctx),5*time.Second);defer stop()
 return s.repo.FinishExecution(finish,e,outcome,response)
}
func(s *Service) Execution(ctx context.Context,actor entity.Operator,id string)(entity.Execution,error){
 if !validID(id){return entity.Execution{},ErrInvalid}
 return s.repo.Execution(ctx,actor,id)
}
