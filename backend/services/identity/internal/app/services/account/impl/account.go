package impl
import(
 "context"
 "errors"
 "strings"
 "time"
 coresearch "grovehub/libs/core/search"
 tx "grovehub/libs/core/database/transaction"
 foundation "grovehub/services/foundation/client/person"
 dto "grovehub/services/identity/internal/app/dto/account"
 "grovehub/services/identity/internal/app/entity"
 mapper "grovehub/services/identity/internal/app/mappers/account"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
 sessionrepo "grovehub/services/identity/internal/app/repository/session"
 "grovehub/services/identity/internal/app/security/principal"
 contract "grovehub/services/identity/internal/app/services/account"
 cred "grovehub/services/identity/internal/app/services/credential"
 invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
 validators "grovehub/services/identity/internal/app/validators/account"
 authhelpers "grovehub/services/identity/internal/app/services/authentication/helpers"
 createop "grovehub/services/identity/internal/app/services/account/operations/create"
 producer "grovehub/services/identity/internal/app/events/producers/account"
 searchop "grovehub/services/identity/internal/app/services/account/operations/search"
)
type Service struct{accounts accountrepo.Repository;tokens tokenrepo.Repository;sessions sessionrepo.Repository;foundation foundation.Client;credential cred.Service;invitation *invitation.Service;mapper mapper.Mapper;validator validators.Validator;unit tx.UnitOfWork;create *createop.Operation;search *searchop.Operation;events producer.ChangedProducer}
func New(a accountrepo.Repository,t tokenrepo.Repository,sessions sessionrepo.Repository,f foundation.Client,c cred.Service,i *invitation.Service,m mapper.Mapper,v validators.Validator,u tx.UnitOfWork,events producer.ChangedProducer)contract.Service{return &Service{accounts:a,tokens:t,sessions:sessions,foundation:f,credential:c,invitation:i,mapper:m,validator:v,unit:u,create:createop.New(a,f,v,i,u,events),search:searchop.New(a,m,v),events:events}}
func normalizeEmail(s string)string{return strings.ToLower(strings.TrimSpace(s))}
func projection(a entity.Account)string{u:="";if a.Username!=nil{u=*a.Username};return strings.ToLower(strings.Join([]string{a.Email,u}," "))}


func(s *Service)Get(ctx context.Context,id string)(dto.Account,error){
 p,err:=principal.RequireAdmin(ctx);if err!=nil{return dto.Account{},err}
 a,err:=s.accounts.FindByID(ctx,id);if err!=nil{return dto.Account{},err}
 if a.TenantID==nil||*a.TenantID!=p.TenantID{return dto.Account{},accountrepo.ErrNotFound}
 return s.mapper.ToResponse(a),nil
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
  before:=a
  a.Version++;a.UpdatedAt=time.Now().UTC()
  if e=s.accounts.Update(inside,a);e!=nil{return e}
  if a.Status==entity.Suspended||a.Status==entity.Deactivated{if _,e=s.sessions.RevokeAll(inside,a.AccountID);e!=nil{return e}}
  if e=s.events.Enqueue(inside,strings.ToUpper(action)+"D",&before,&a);e!=nil{return e}
  out=a;return nil
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
  if e=s.events.Enqueue(inside,"ACTIVATED",nil,&a);e!=nil{return e}
  out=dto.SetupResult{AccountID:id,PersonGID:a.PersonGID,Username:name,Status:string(a.Status),SetupCompleted:true}
  return nil
 })
 return out,err
}

func(s *Service)Create(ctx context.Context,input dto.Create,key string)(dto.CreateResult,error){return s.create.Execute(ctx,input,key)}

func(s *Service)Search(ctx context.Context,request coresearch.Request)(coresearch.Result[dto.Account],error){return s.search.Execute(ctx,request)}
