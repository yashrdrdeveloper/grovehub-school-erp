package impl
import(
 "context"
 "errors"
 "fmt"
 "strings"
 "time"
 "github.com/google/uuid"
 coreid "grovehub/libs/core/id"
 coresearch "grovehub/libs/core/search"
 tx "grovehub/libs/core/database/transaction"
 coreidempotency "grovehub/libs/core/idempotency"
 headercontext "grovehub/libs/core/headercontext"
 foundation "grovehub/services/foundation/client/person"
 dto "grovehub/services/identity/internal/app/dto/account"
 "grovehub/services/identity/internal/app/entity"
 mapper "grovehub/services/identity/internal/app/mappers/account"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
 "grovehub/services/identity/internal/app/security/principal"
 contract "grovehub/services/identity/internal/app/services/account"
 cred "grovehub/services/identity/internal/app/services/credential"
 invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
 validators "grovehub/services/identity/internal/app/validators/account"
 authhelpers "grovehub/services/identity/internal/app/services/authentication/helpers"
)
type Service struct{accounts accountrepo.Repository;tokens tokenrepo.Repository;foundation foundation.Client;credential cred.Service;invitation *invitation.Service;mapper mapper.Mapper;validator validators.Validator;unit tx.UnitOfWork}
func New(a accountrepo.Repository,t tokenrepo.Repository,f foundation.Client,c cred.Service,i *invitation.Service,m mapper.Mapper,v validators.Validator,u tx.UnitOfWork)contract.Service{return &Service{accounts:a,tokens:t,foundation:f,credential:c,invitation:i,mapper:m,validator:v,unit:u}}
func normalizeEmail(s string)string{return strings.ToLower(strings.TrimSpace(s))}
func projection(a entity.Account)string{u:="";if a.Username!=nil{u=*a.Username};return strings.ToLower(strings.Join([]string{a.Email,u}," "))}
func(s *Service)Create(ctx context.Context,input dto.Create,idempotencyKey string)(dto.CreateResult,error){
 p,err:=principal.RequireAdmin(ctx);if err!=nil{return dto.CreateResult{},err}
 if s.invitation==nil{return dto.CreateResult{},errors.New("invitation sender unavailable")}
 input,err=s.validator.NormalizeCreate(input);if err!=nil{return dto.CreateResult{},err}
 if _,err=coreidempotency.ValidateKey(idempotencyKey);err!=nil{return dto.CreateResult{},err}
 fingerprint,err:=coreidempotency.Fingerprint(input);if err!=nil{return dto.CreateResult{},err}
 var result dto.CreateResult
 err=s.unit.Run(ctx,func(inside context.Context)error{
  current,lookup:=s.accounts.FindByIdempotencyKey(inside,idempotencyKey)
  if lookup==nil{
   if current.RegistrationRequestHash==nil||*current.RegistrationRequestHash!=fingerprint{return errors.New("idempotency key conflict")}
   if current.TenantID==nil||*current.TenantID!=p.TenantID{return errors.New("unauthorized tenant")}
   result=convertCreate(current);return nil
  }
  if !errors.Is(lookup,accountrepo.ErrNotFound){return lookup}
  person,personErr:=s.resolvePerson(inside,input);if personErr!=nil{return personErr}
  _,existingErr:=s.accounts.FindByPersonGID(inside,person.GID)
  if existingErr==nil{return accountrepo.ErrConflict}
  if !errors.Is(existingErr,accountrepo.ErrNotFound){return existingErr}
  id,e:=coreid.NewUUIDv7();if e!=nil{return e}
  now:=time.Now().UTC()
  a:=entity.Account{AccountID:id,PersonGID:person.GID,TenantID:&p.TenantID,Email:input.Email,Status:entity.PendingSetup,Version:1,CreatedAt:now,UpdatedAt:now,RegistrationIdempotencyKey:&idempotencyKey,RegistrationRequestHash:&fingerprint}
  if input.PhoneNumber!=""{a.Phone=&input.PhoneNumber}
  a.SearchDetails=projection(a)
  if e=s.accounts.Create(inside,a);e!=nil{return e}
  if e=s.invitation.Queue(inside,a,entity.AccountSetup,input.Person.FirstName);e!=nil{return e}
  result=convertCreate(a);return nil
 })
 return result,err
}
func convertCreate(a entity.Account)dto.CreateResult{
 r:=dto.CreateResult{AccountID:a.AccountID,PersonGID:a.PersonGID,Email:a.Email,Status:string(a.Status)}
 r.Invitation.Status="QUEUED";return r
}
func(s *Service)resolvePerson(ctx context.Context,input dto.Create)(foundation.Person,error){
 query:=coresearch.Request{Filters:[]coresearch.Filter{{Field:"contacts.email",Operator:coresearch.OperatorEquals,Values:[]string{input.Email}}},PageInfo:coresearch.PageRequest{Rows:2}}
 result,err:=s.foundation.SearchPersons(ctx,foundation.SearchPersonsRequest{RequestBody:query})
 if err!=nil{return foundation.Person{},err}
 if len(result.Items)>1{return foundation.Person{},errors.New("ambiguous Foundation person")}
 if len(result.Items)==1{
  for _,e:=range result.Items[0].Contacts.Emails{if normalizeEmail(e.Value)==input.Email{return result.Items[0],nil}}
  return foundation.Person{},errors.New("Foundation email mismatch")
 }
 middle:="";if input.Person.MiddleName!=nil{middle=*input.Person.MiddleName}
 newPerson:=foundation.CreatePerson{PersonNames:foundation.PersonNames{FirstName:input.Person.FirstName,MiddleName:middle,LastName:input.Person.LastName},Contacts:foundation.Contacts{Emails:[]foundation.ContactValue{{Value:input.Email,IsPrimary:true}}}}
 if input.PhoneNumber!=""{newPerson.Contacts.Phones=[]foundation.ContactValue{{Value:input.PhoneNumber,IsPrimary:true}}}
 created,_,err:=s.foundation.CreatePerson(headercontext.WithoutIdempotencyKey(ctx),foundation.CreatePersonRequest{RequestBody:newPerson})
 return created,err
}
func(s *Service)Get(ctx context.Context,id string)(dto.Account,error){
 p,err:=principal.RequireAdmin(ctx);if err!=nil{return dto.Account{},err}
 a,err:=s.accounts.FindByID(ctx,id);if err!=nil{return dto.Account{},err}
 if a.TenantID==nil||*a.TenantID!=p.TenantID{return dto.Account{},accountrepo.ErrNotFound}
 return s.mapper.ToResponse(a),nil
}
func(s *Service)Search(ctx context.Context,request coresearch.Request)(coresearch.Result[dto.Account],error){
 p,err:=principal.RequireAdmin(ctx);if err!=nil{return coresearch.Result[dto.Account]{},err}
 q,err:=s.validator.Search(request);if err!=nil{return coresearch.Result[dto.Account]{},err}
 dbResult,err:=s.accounts.Search(ctx,p.TenantID,q);if err!=nil{return coresearch.Result[dto.Account]{},err}
 out:=coresearch.Result[dto.Account]{Items:make([]dto.Account,0,len(dbResult.Items))}
 for _,a:=range dbResult.Items{out.Items=append(out.Items,s.mapper.ToResponse(a))}
 out.PageInfo=coresearch.PageResponse{Rows:len(out.Items),HasNext:dbResult.HasNext}
 if dbResult.HasNext&&len(dbResult.Items)>0{
  last:=dbResult.Items[len(dbResult.Items)-1];value:=last.CreatedAt.UTC().Format(time.RFC3339Nano)
  switch q.Sort.Field{case "accountId":value=last.AccountID;case "updatedAt":value=last.UpdatedAt.UTC().Format(time.RFC3339Nano);case "status":value=string(last.Status);case "username":if last.Username!=nil{value=*last.Username}else{value=""}}
  out.PageInfo.NextCursor,err=coresearch.EncodeCursor(coresearch.CursorPayload{Field:q.Sort.Field,Direction:q.Sort.Direction,Value:value,CursorId:last.AccountID})
 }
 return out,err
}
func(s *Service)Update(ctx context.Context,id string,request dto.Update)(dto.Account,error){
 p,err:=principal.RequireAdmin(ctx);if err!=nil{return dto.Account{},err}
 if request.Username==nil{return dto.Account{},errors.New("username must be supplied")}
 name,err:=s.validator.NormalizeUsername(*request.Username);if err!=nil{return dto.Account{},err}
 var updated entity.Account
 err=s.unit.Run(ctx,func(inside context.Context)error{
  a,e:=s.accounts.FindByID(inside,id);if e!=nil{return e}
  if a.TenantID==nil||*a.TenantID!=p.TenantID{return accountrepo.ErrNotFound}
  a.Username=&name;a.Version++;a.UpdatedAt=time.Now().UTC();a.SearchDetails=projection(a)
  if e=s.accounts.Update(inside,a);e!=nil{return e};updated=a;return nil
 })
 return s.mapper.ToResponse(updated),err
}
func(s *Service)Transition(ctx context.Context,id,action string)(dto.Account,error){
 p,err:=principal.RequireAdmin(ctx);if err!=nil{return dto.Account{},err}
 var out entity.Account
 err=s.unit.Run(ctx,func(inside context.Context)error{
  a,e:=s.accounts.FindByID(inside,id);if e!=nil{return e}
  if a.TenantID==nil||*a.TenantID!=p.TenantID{return accountrepo.ErrNotFound}
  switch action{
  case "suspend":if a.Status!=entity.Active{return errors.New("invalid status transition")};a.Status=entity.Suspended
  case "activate":if a.Status!=entity.Suspended{return errors.New("invalid status transition")};a.Status=entity.Active
  case "deactivate":if a.Status==entity.Deactivated{return errors.New("already deactivated")};a.Status=entity.Deactivated
  default:return errors.New("unknown account transition")
  }
  a.Version++;a.UpdatedAt=time.Now().UTC()
  if e=s.accounts.Update(inside,a);e!=nil{return e};out=a;return nil
 })
 return s.mapper.ToResponse(out),err
}
func(s *Service)Resend(ctx context.Context,id string)error{
 p,err:=principal.RequireAdmin(ctx);if err!=nil{return err}
 return s.unit.Run(ctx,func(inside context.Context)error{
  a,e:=s.accounts.FindByID(inside,id);if e!=nil{return e}
  if a.TenantID==nil||*a.TenantID!=p.TenantID{return accountrepo.ErrNotFound}
  if a.Status!=entity.PendingSetup{return errors.New("account is not pending setup")}
  return s.invitation.Queue(inside,a,entity.AccountSetup,"there")
 })
}
func(s *Service)Setup(ctx context.Context,grant string,input dto.Setup)(dto.SetupResult,error){
 name,err:=s.validator.NormalizeUsername(input.Username);if err!=nil{return dto.SetupResult{},err}
 digest,err:=authhelpers.Digest(grant,32);if err!=nil{return dto.SetupResult{},err}
 var out dto.SetupResult
 err=s.unit.Run(ctx,func(inside context.Context)error{
  id,e:=s.tokens.Consume(inside,digest,entity.SetupGrant,time.Now().UTC());if e!=nil{return e}
  a,e:=s.accounts.FindByID(inside,id);if e!=nil{return e}
  if a.Status!=entity.PendingSetup{return errors.New("account not pending")}
  a.Username=&name;a.Status=entity.Active;a.Version++;a.UpdatedAt=time.Now().UTC();a.SearchDetails=projection(a)
  if e=s.credential.Set(inside,id,input.Password);e!=nil{return e}
  if e=s.accounts.Update(inside,a);e!=nil{return e}
  out=dto.SetupResult{AccountID:id,PersonGID:a.PersonGID,Username:name,Status:string(a.Status),SetupCompleted:true}
  return nil
 })
 return out,err
}
var _=uuid.NewString
var _=fmt.Sprint
