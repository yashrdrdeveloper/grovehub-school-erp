package impl
import(
 "context"
 "errors"
 "strings"
 "time"
 "github.com/google/uuid"
 coreid "grovehub/libs/core/id"
 tx "grovehub/libs/core/database/transaction"
 accountdto "grovehub/services/identity/internal/app/dto/account"
 authdto "grovehub/services/identity/internal/app/dto/authentication"
 "grovehub/services/identity/internal/app/entity"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
 sessionrepo "grovehub/services/identity/internal/app/repository/session"
 accountmapper "grovehub/services/identity/internal/app/mappers/account"
 "grovehub/services/identity/internal/app/security/principal"
 contract "grovehub/services/identity/internal/app/services/authentication"
 secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
 credential "grovehub/services/identity/internal/app/services/credential"
 invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
)
type Service struct{accounts accountrepo.Repository;sessions sessionrepo.Repository;tokens tokenrepo.Repository;credential credential.Service;invitation *invitation.Service;mapper accountmapper.Mapper;unit tx.UnitOfWork}
func New(a accountrepo.Repository,s sessionrepo.Repository,t tokenrepo.Repository,c credential.Service,i *invitation.Service,m accountmapper.Mapper,u tx.UnitOfWork)contract.Service{return &Service{accounts:a,sessions:s,tokens:t,credential:c,invitation:i,mapper:m,unit:u}}
func(s *Service)Login(ctx context.Context,input authdto.Login,userAgent string)(contract.SessionResult,error){
 identifier:=strings.ToLower(strings.TrimSpace(input.Identifier))
 if identifier==""||input.Password==""{return contract.SessionResult{},errors.New("invalid credentials")}
 var a entity.Account
 var err error
 if strings.Contains(identifier,"@"){a,err=s.accounts.FindByEmail(ctx,identifier)}else{a,err=s.accounts.FindByUsername(ctx,identifier)}
 if err!=nil||a.Status!=entity.Active{return contract.SessionResult{},errors.New("invalid credentials")}
 ok,err:=s.credential.Verify(ctx,a.AccountID,input.Password);if err!=nil||!ok{return contract.SessionResult{},errors.New("invalid credentials")}
 access,err:=secrets.NewSecret(32);if err!=nil{return contract.SessionResult{},err}
 refresh,err:=secrets.NewSecret(48);if err!=nil{return contract.SessionResult{},err}
 ah,_:=secrets.Digest(access,32);rh,_:=secrets.Digest(refresh,48)
 id,err:=coreid.NewUUIDv7();if err!=nil{return contract.SessionResult{},err}
 now:=time.Now().UTC()
 session:=entity.Session{SessionID:id,AccountID:a.AccountID,AccessHash:ah,RefreshHash:rh,AccessExpiresAt:now.Add(15*time.Minute),RefreshExpiresAt:now.Add(30*24*time.Hour),ReauthenticatedAt:&now,UserAgent:truncate(userAgent,255),CreatedAt:now,UpdatedAt:now}
 if err=s.sessions.Create(ctx,session);err!=nil{return contract.SessionResult{},err}
 return contract.SessionResult{Data:authdto.LoginResult{AccountID:a.AccountID,PersonGID:a.PersonGID,SessionID:id,ExpiresIn:900,RequiresSetup:false},Cookies:contract.Cookies{Access:access,Refresh:refresh}},nil
}
func truncate(s string,n int)string{if len(s)>n{return s[:n]};return s}
func(s *Service)Refresh(ctx context.Context,old string)(contract.RefreshResult,error){
 hash,err:=secrets.Digest(old,48);if err!=nil{return contract.RefreshResult{},errors.New("invalid refresh token")}
 existing,err:=s.sessions.FindRefresh(ctx,hash,time.Now().UTC());if err!=nil{return contract.RefreshResult{},errors.New("invalid refresh token")}
 account,err:=s.accounts.FindByID(ctx,existing.AccountID);if err!=nil||account.Status!=entity.Active{return contract.RefreshResult{},errors.New("invalid refresh token")}
 access,err:=secrets.NewSecret(32);if err!=nil{return contract.RefreshResult{},err}
 refresh,err:=secrets.NewSecret(48);if err!=nil{return contract.RefreshResult{},err}
 now:=time.Now().UTC();ah,_:=secrets.Digest(access,32);rh,_:=secrets.Digest(refresh,48)
 updated:=existing;updated.AccessHash=ah;updated.RefreshHash=rh;updated.AccessExpiresAt=now.Add(15*time.Minute);updated.UpdatedAt=now
 if err=s.sessions.Rotate(ctx,updated,hash);err!=nil{return contract.RefreshResult{},errors.New("invalid refresh token")}
 return contract.RefreshResult{Data:authdto.RefreshResult{SessionID:existing.SessionID,ExpiresIn:900,Refreshed:true},Cookies:contract.Cookies{Access:access,Refresh:refresh}},nil
}
func(s *Service)Forgot(ctx context.Context,email string)error{
 email=strings.ToLower(strings.TrimSpace(email));if !strings.Contains(email,"@"){return nil}
 a,err:=s.accounts.FindByEmail(ctx,email);if err!=nil||a.Status!=entity.Active{return nil}
 if s.invitation==nil{return errors.New("mail unavailable")}
 return s.unit.Run(ctx,func(inside context.Context)error{return s.invitation.Queue(inside,a,entity.PasswordReset,"there")})
}
func(s *Service)Exchange(ctx context.Context,raw,purpose string)(contract.Grant,error){
 var original,grant entity.TokenPurpose
 switch purpose{case "recovery":original=entity.PasswordReset;grant=entity.ResetGrant;case "invitation":original=entity.AccountSetup;grant=entity.SetupGrant;default:return contract.Grant{},errors.New("invalid proof purpose")}
 digest,err:=secrets.Digest(raw,32);if err!=nil{return contract.Grant{},errors.New("invalid proof")}
 newSecret,err:=secrets.NewSecret(32);if err!=nil{return contract.Grant{},err}
 grantHash,_:=secrets.Digest(newSecret,32)
 err=s.unit.Run(ctx,func(inside context.Context)error{
  id,e:=s.tokens.Consume(inside,digest,original,time.Now().UTC());if e!=nil{return e}
  a,e:=s.accounts.FindByID(inside,id);if e!=nil{return e}
  if (original==entity.AccountSetup&&a.Status!=entity.PendingSetup)||(original==entity.PasswordReset&&a.Status!=entity.Active){return errors.New("invalid account state")}
  now:=time.Now().UTC()
  return s.tokens.Issue(inside,entity.PasswordToken{TokenHash:grantHash,AccountID:id,Purpose:grant,ExpiresAt:now.Add(10*time.Minute),CreatedAt:now})
 })
 if err!=nil{return contract.Grant{},errors.New("invalid or expired proof")}
 return contract.Grant{Secret:newSecret,Data:authdto.GrantResult{Verified:true,Purpose:string(original),ExpiresIn:600}},nil
}
func(s *Service)ResetPassword(ctx context.Context,grant,password string)error{
 hash,err:=secrets.Digest(grant,32);if err!=nil{return errors.New("invalid reset grant")}
 return s.unit.Run(ctx,func(inside context.Context)error{
  id,e:=s.tokens.Consume(inside,hash,entity.ResetGrant,time.Now().UTC());if e!=nil{return e}
  a,e:=s.accounts.FindByID(inside,id);if e!=nil{return e};if a.Status!=entity.Active{return errors.New("invalid account status")}
  if e=s.credential.Replace(inside,id,password);e!=nil{return e}
  _,e=s.sessions.RevokeAll(inside,id);return e
 })
}
func(s *Service)Reauthenticate(ctx context.Context,password string)error{
 p,err:=principal.Require(ctx);if err!=nil{return err}
 valid,err:=s.credential.Verify(ctx,p.AccountID,password);if err!=nil||!valid{return errors.New("invalid credentials")}
 return s.sessions.Reauthenticate(ctx,p.SessionID,time.Now().UTC())
}
func(s *Service)ChangePassword(ctx context.Context,password string)error{
 p,err:=principal.Require(ctx);if err!=nil{return err}
 if err=principal.Recent(ctx,5*time.Minute);err!=nil{return err}
 return s.unit.Run(ctx,func(inside context.Context)error{
  if e:=s.credential.Replace(inside,p.AccountID,password);e!=nil{return e}
  _,e:=s.sessions.RevokeAll(inside,p.AccountID);return e
 })
}
func(s *Service)Logout(ctx context.Context)error{p,err:=principal.Require(ctx);if err!=nil{return err};return s.sessions.Revoke(ctx,p.AccountID,p.SessionID)}
func(s *Service)LogoutAll(ctx context.Context)(int64,error){
 p,err:=principal.Require(ctx);if err!=nil{return 0,err}
 if err=principal.Recent(ctx,5*time.Minute);err!=nil{return 0,err}
 return s.sessions.RevokeAll(ctx,p.AccountID)
}
func(s *Service)ListSessions(ctx context.Context)([]map[string]any,error){
 p,err:=principal.Require(ctx);if err!=nil{return nil,err}
 sessions,err:=s.sessions.List(ctx,p.AccountID);if err!=nil{return nil,err}
 out:=make([]map[string]any,0,len(sessions))
 for _,item:=range sessions{out=append(out,map[string]any{"sessionId":item.SessionID,"current":item.SessionID==p.SessionID,"createdAt":item.CreatedAt,"lastActiveAt":item.LastActiveAt,"expiresAt":item.RefreshExpiresAt})}
 return out,nil
}
func(s *Service)RevokeSession(ctx context.Context,id string)error{
 p,err:=principal.Require(ctx);if err!=nil{return err}
 if _,e:=uuid.Parse(id);e!=nil{return errors.New("invalid session ID")}
 return s.sessions.Revoke(ctx,p.AccountID,id)
}
func(s *Service)Me(ctx context.Context)(accountdto.Account,error){
 p,err:=principal.Require(ctx);if err!=nil{return accountdto.Account{},err}
 a,err:=s.accounts.FindByID(ctx,p.AccountID);if err!=nil{return accountdto.Account{},err}
 return s.mapper.ToResponse(a),nil
}
