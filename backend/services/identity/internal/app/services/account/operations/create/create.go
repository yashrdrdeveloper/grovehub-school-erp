package create
import(
 "context"
 "errors"
 "strings"
 "time"
 coreid "grovehub/libs/core/id"
 coreidempotency "grovehub/libs/core/idempotency"
 headercontext "grovehub/libs/core/headercontext"
 tx "grovehub/libs/core/database/transaction"
 corepipeline "grovehub/libs/core/pipeline/service"
 coresearch "grovehub/libs/core/search"
 foundation "grovehub/services/foundation/client/person"
 dto "grovehub/services/identity/internal/app/dto/account"
 "grovehub/services/identity/internal/app/entity"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 "grovehub/services/identity/internal/app/security/principal"
 invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
 validators "grovehub/services/identity/internal/app/validators/account"
 "grovehub/services/identity/internal/app/services/account/operations/create/flows"
)
type Operation struct{
 accounts accountrepo.Repository
 foundation foundation.Client
 validator validators.Validator
 invitation *invitation.Service
 unit tx.UnitOfWork
}
func New(a accountrepo.Repository,f foundation.Client,v validators.Validator,i *invitation.Service,u tx.UnitOfWork)*Operation{return &Operation{accounts:a,foundation:f,validator:v,invitation:i,unit:u}}
// @ChainOfResponsibility
func(o *Operation)Execute(ctx context.Context,request dto.Create,key string)(dto.CreateResult,error){
 flow:=&flows.Create{Request:request,IdempotencyKey:key}
 wrap:=func(step func(context.Context,*flows.Create)error)corepipeline.Handler[flows.Create,dto.CreateResult]{
  return func(ctx context.Context,f *flows.Create)corepipeline.Outcome[dto.CreateResult]{
   if err:=step(ctx,f);err!=nil{return corepipeline.Fail[dto.CreateResult](err)}
   return corepipeline.Next[dto.CreateResult]()
  }
 }
 var result dto.CreateResult
 err:=o.unit.Run(ctx,func(inside context.Context)error{
  var runErr error
  result,_,runErr=corepipeline.Run(inside,flow,
   wrap(o.authorize),
   wrap(o.validate),
   func(ctx context.Context,f *flows.Create)corepipeline.Outcome[dto.CreateResult]{
    if err:=o.replay(ctx,f);err!=nil{return corepipeline.Fail[dto.CreateResult](err)}
    if f.Done{return corepipeline.Return(f.Response,false)}
    return corepipeline.Next[dto.CreateResult]()
   },
   wrap(o.resolve),
   func(ctx context.Context,f *flows.Create)corepipeline.Outcome[dto.CreateResult]{
    if err:=o.persist(ctx,f);err!=nil{return corepipeline.Fail[dto.CreateResult](err)}
    return corepipeline.Return(f.Response,true)
   },
  )
  return runErr
 })
 return result,err
}
func(o *Operation)authorize(ctx context.Context,f *flows.Create)error{
 p,err:=principal.RequireAdmin(ctx);if err!=nil{return err}
 if o.invitation==nil{return errors.New("invitation service unavailable")}
 f.Actor=p;return nil
}
func(o *Operation)validate(ctx context.Context,f *flows.Create)error{
 r,err:=o.validator.NormalizeCreate(f.Request);if err!=nil{return err}
 f.Request=r
 if _,err=coreidempotency.ValidateKey(f.IdempotencyKey);err!=nil{return err}
 f.RequestHash,err=coreidempotency.Fingerprint(f.Request);return err
}
func(o *Operation)replay(ctx context.Context,f *flows.Create)error{
 a,err:=o.accounts.FindByIdempotencyKey(ctx,f.IdempotencyKey)
 if errors.Is(err,accountrepo.ErrNotFound){return nil}
 if err!=nil{return err}
 if a.RegistrationRequestHash==nil||*a.RegistrationRequestHash!=f.RequestHash{return errors.New("idempotency key conflict")}
 if a.TenantID==nil||*a.TenantID!=f.Actor.TenantID{return errors.New("forbidden tenant")}
 f.Account=a;f.Response=mapCreated(a);f.Done=true;return nil
}
func(o *Operation)resolve(ctx context.Context,f *flows.Create)error{
 r:=f.Request
 search:=coresearch.Request{Filters:[]coresearch.Filter{{Field:"contacts.email",Operator:coresearch.OperatorEquals,Values:[]string{r.Email}}},PageInfo:coresearch.PageRequest{Rows:2}}
 result,err:=o.foundation.SearchPersons(ctx,foundation.SearchPersonsRequest{RequestBody:search});if err!=nil{return err}
 if len(result.Items)>1{return errors.New("ambiguous person")}
 if len(result.Items)==1{
  found:=false;for _,email:=range result.Items[0].Contacts.Emails{if strings.EqualFold(email.Value,r.Email){found=true;break}}
  if !found{return errors.New("Foundation contact mismatch")}
  f.Person=result.Items[0];return nil
 }
 middle:="";if r.Person.MiddleName!=nil{middle=*r.Person.MiddleName}
 create:=foundation.CreatePerson{PersonNames:foundation.PersonNames{FirstName:r.Person.FirstName,MiddleName:middle,LastName:r.Person.LastName},Contacts:foundation.Contacts{Emails:[]foundation.ContactValue{{Value:r.Email,IsPrimary:true}}}}
 if r.PhoneNumber!=""{create.Contacts.Phones=[]foundation.ContactValue{{Value:r.PhoneNumber,IsPrimary:true}}}
 f.Person,_,err=o.foundation.CreatePerson(headercontext.WithoutIdempotencyKey(ctx),foundation.CreatePersonRequest{RequestBody:create});return err
}
func(o *Operation)persist(ctx context.Context,f *flows.Create)error{
 if _,err:=o.accounts.FindByPersonGID(ctx,f.Person.GID);err==nil{return accountrepo.ErrConflict}else if !errors.Is(err,accountrepo.ErrNotFound){return err}
 id,err:=coreid.NewUUIDv7();if err!=nil{return err}
 now:=time.Now().UTC()
 a:=entity.Account{AccountID:id,PersonGID:f.Person.GID,TenantID:&f.Actor.TenantID,Email:f.Request.Email,Status:entity.PendingSetup,Version:1,CreatedAt:now,UpdatedAt:now,RegistrationIdempotencyKey:&f.IdempotencyKey,RegistrationRequestHash:&f.RequestHash}
 if f.Request.PhoneNumber!=""{a.Phone=&f.Request.PhoneNumber}
 a.SearchDetails=strings.ToLower(a.Email)
 if err=o.accounts.Create(ctx,a);err!=nil{return err}
 if err=o.invitation.Queue(ctx,a,entity.AccountSetup,f.Request.Person.FirstName);err!=nil{return err}
 f.Account=a;f.Response=mapCreated(a);return nil
}
func mapCreated(a entity.Account)dto.CreateResult{
 out:=dto.CreateResult{AccountID:a.AccountID,PersonGID:a.PersonGID,Email:a.Email,Status:string(a.Status)}
 out.Invitation.Status="QUEUED";return out
}
