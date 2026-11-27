package workflow

import (
 "context"
 "database/sql"
 "encoding/json"
 "errors"
 "regexp"
 "strings"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
 repo "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/repository/workflow"
)

var ErrInvalid=errors.New("invalid workflow request")
var ErrForbidden=repo.ErrForbidden
var ErrConflict=repo.ErrConflict

var codePattern=regexp.MustCompile("^[a-z][a-z0-9_-]{1,63}$")
var fieldPattern=regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]{0,63}$")
var endpointPattern=regexp.MustCompile("^/[a-z][a-z0-9-]*/restricted(/[a-zA-Z0-9_.-]+|/\\$\\{[a-zA-Z][a-zA-Z0-9_]*\\})+$")
var placeholderPattern=regexp.MustCompile("\\$\\{([a-zA-Z][a-zA-Z0-9_]*)\\}")
type Service struct{repo *repo.Repository}
func New(db *sql.DB)*Service{return &Service{repo.New(db)}}

type EnvironmentInput struct{Code string `json:"code"`;Name string `json:"name"`}
type TemplateInput struct {
 Code string `json:"code"`;Module string `json:"module"`;Method string `json:"method"`;Endpoint string `json:"endpoint"`;Scope string `json:"scope"`
 Headers map[string]string `json:"headers"`;BodySchema map[string]entity.Field `json:"bodySchema"`
}
type SubmitInput struct{
 TemplateID string `json:"templateId"`;EnvironmentID string `json:"environmentId"`;SchoolID *string `json:"schoolId"`;Input json.RawMessage `json:"input"`;Reason string `json:"reason"`
}
type DecisionInput struct{Version int64 `json:"version"`;Comment string `json:"comment"`}
func validID(s string)bool{v,err:=uuid.Parse(s);return err==nil&&v!=uuid.Nil&&v.String()==s}
func(s *Service) CreateEnvironment(ctx context.Context,actor entity.Operator,in EnvironmentInput)(entity.Environment,error){
 if actor.Persona!=entity.L2{return entity.Environment{},ErrForbidden}
 in.Name=strings.TrimSpace(in.Name)
 if !codePattern.MatchString(in.Code)||in.Name==""||len(in.Name)>100{return entity.Environment{},ErrInvalid}
 e:=entity.Environment{ID:uuid.NewString(),Code:in.Code,Name:in.Name};return e,s.repo.CreateEnvironment(ctx,actor.ID,e)
}
func(s *Service) Environments(ctx context.Context)([]entity.Environment,error){return s.repo.Environments(ctx)}
func(s *Service) Templates(ctx context.Context)([]entity.Template,error){return s.repo.Templates(ctx)}
func validateTemplate(in TemplateInput)bool{
 if !codePattern.MatchString(in.Code)||!codePattern.MatchString(in.Module)||len(in.Endpoint)>512||!endpointPattern.MatchString(in.Endpoint)||!strings.HasPrefix(in.Endpoint,"/"+in.Module+"/restricted/"){return false}
 for _,segment:=range strings.Split(in.Endpoint,"/"){if segment=="."||segment==".."{return false}}
 switch in.Method{case "GET","POST","PATCH","DELETE":default:return false}
 if in.Scope!="GLOBAL"&&in.Scope!="SCHOOL"{return false}
 if len(in.Headers)>3||len(in.BodySchema)>32{return false}
 for key,field:=range in.BodySchema{
  if !fieldPattern.MatchString(key){return false}
  switch field.Type{case "string","integer","boolean":default:return false}
 }
 for name,value:=range in.Headers{
  switch name{case "Content-Type":if value!="application/json"{return false};case "Idempotency-Key":if len(value)>128||!placeholderPattern.MatchString(value){return false};default:return false}
  if strings.ContainsAny(value,"\r\n"){return false}
  if name=="Idempotency-Key"{m:=placeholderPattern.FindStringSubmatch(value);if m[0]!=value||!in.BodySchema[m[1]].Required||in.BodySchema[m[1]].Type!="string"{return false}}
 }
 for _,m:=range placeholderPattern.FindAllStringSubmatch(in.Endpoint,-1){if !in.BodySchema[m[1]].Required||in.BodySchema[m[1]].Type!="string"{return false}}
 return true
}
func(s *Service) CreateTemplate(ctx context.Context,actor entity.Operator,in TemplateInput)(entity.Template,error){
 if actor.Persona!=entity.L2{return entity.Template{},ErrForbidden}
 if !validateTemplate(in){return entity.Template{},ErrInvalid}
 if in.Headers==nil{in.Headers=map[string]string{}};if in.BodySchema==nil{in.BodySchema=map[string]entity.Field{}}
 t:=entity.Template{ID:uuid.NewString(),Code:in.Code,Module:in.Module,Method:in.Method,Endpoint:in.Endpoint,Scope:in.Scope,Headers:in.Headers,BodySchema:in.BodySchema}
 err:=s.repo.CreateTemplate(ctx,actor.ID,&t);return t,err
}
func validInput(raw json.RawMessage,schema map[string]entity.Field)bool{
 if len(raw)==0||len(raw)>16384{return false}
 var values map[string]json.RawMessage
 if json.Unmarshal(raw,&values)!=nil||values==nil{return false}
 for key,value:=range values{
  field,ok:=schema[key];if !ok||string(value)=="null"{return false}
  switch field.Type{
  case "string":var v string;if json.Unmarshal(value,&v)!=nil||len(v)>4096{return false}
  case "integer":var v int64;if json.Unmarshal(value,&v)!=nil{return false}
  case "boolean":var v bool;if json.Unmarshal(value,&v)!=nil{return false}
  default:return false
  }
 }
 for key,field:=range schema{if field.Required{if _,ok:=values[key];!ok{return false}}}
 return true
}
func(s *Service) Submit(ctx context.Context,actor entity.Operator,in SubmitInput)(entity.Approval,error){
 if actor.Persona!=entity.L1&&actor.Persona!=entity.L2{return entity.Approval{},ErrForbidden}
 in.Reason=strings.TrimSpace(in.Reason)
 if !validID(in.TemplateID)||!validID(in.EnvironmentID)||in.Reason==""||len(in.Reason)>1000{return entity.Approval{},ErrInvalid}
 t,err:=s.repo.Template(ctx,in.TemplateID);if err!=nil{return entity.Approval{},err}
 if (t.Scope=="GLOBAL"&&in.SchoolID!=nil)||(t.Scope=="SCHOOL"&&(in.SchoolID==nil||!validID(*in.SchoolID)))||!validInput(in.Input,t.BodySchema){return entity.Approval{},ErrInvalid}
 return s.repo.Submit(ctx,entity.Approval{ID:uuid.NewString(),TemplateID:t.ID,EnvironmentID:in.EnvironmentID,SchoolID:in.SchoolID,Input:in.Input,Reason:in.Reason,RequesterID:actor.ID},actor.Persona==entity.L2)
}
func(s *Service) List(ctx context.Context,actor entity.Operator,inbox bool)([]entity.Approval,error){return s.repo.List(ctx,actor,inbox)}
func(s *Service) Get(ctx context.Context,actor entity.Operator,id string)(entity.Approval,error){if !validID(id){return entity.Approval{},ErrInvalid};return s.repo.Get(ctx,actor,id)}
func(s *Service) Decide(ctx context.Context,actor entity.Operator,id,status string,in DecisionInput)(entity.Approval,error){
 in.Comment=strings.TrimSpace(in.Comment)
 if !validID(id)||in.Version<1||len(in.Comment)>1000||(status!="APPROVED"&&status!="REJECTED"&&status!="CANCELLED")||(status=="REJECTED"&&in.Comment==""){return entity.Approval{},ErrInvalid}
 return s.repo.Decide(ctx,actor,id,in.Version,status,in.Comment)
}
