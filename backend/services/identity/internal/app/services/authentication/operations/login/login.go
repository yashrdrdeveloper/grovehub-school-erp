package login
import(
 "context"
 "errors"
 "strings"
 "time"
 coreid "grovehub/libs/core/id"
 dto "grovehub/services/identity/internal/app/dto/authentication"
 "grovehub/services/identity/internal/app/entity"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 sessionrepo "grovehub/services/identity/internal/app/repository/session"
 auth "grovehub/services/identity/internal/app/services/authentication"
 helpers "grovehub/services/identity/internal/app/services/authentication/helpers"
 credential "grovehub/services/identity/internal/app/services/credential"
)
type Operation struct{accounts accountrepo.Repository;sessions sessionrepo.Repository;credential credential.Service}
func New(a accountrepo.Repository,s sessionrepo.Repository,c credential.Service)*Operation{return &Operation{accounts:a,sessions:s,credential:c}}
func(o *Operation)Execute(ctx context.Context,input dto.Login,userAgent string)(auth.SessionResult,error){
 identifier:=strings.ToLower(strings.TrimSpace(input.Identifier))
 if identifier==""||input.Password==""{return auth.SessionResult{},errors.New("invalid credentials")}
 if !strings.Contains(identifier,"@"){return auth.SessionResult{},errors.New("invalid credentials")}
 account,err:=o.accounts.FindByEmail(ctx,identifier)
 if err!=nil||account.Status!=entity.Active{return auth.SessionResult{},errors.New("invalid credentials")}
 ok,err:=o.credential.Verify(ctx,account.AccountID,input.Password);if err!=nil||!ok{return auth.SessionResult{},errors.New("invalid credentials")}
 access,err:=helpers.NewSecret(32);if err!=nil{return auth.SessionResult{},err}
 refresh,err:=helpers.NewSecret(48);if err!=nil{return auth.SessionResult{},err}
 accessHash,_:=helpers.Digest(access,32);refreshHash,_:=helpers.Digest(refresh,48)
 id,err:=coreid.NewUUIDv7();if err!=nil{return auth.SessionResult{},err}
 now:=time.Now().UTC()
 if len(userAgent)>255{userAgent=userAgent[:255]}
 session:=entity.Session{SessionID:id,AccountID:account.AccountID,AccessHash:accessHash,RefreshHash:refreshHash,AccessExpiresAt:now.Add(15*time.Minute),RefreshExpiresAt:now.Add(30*24*time.Hour),ReauthenticatedAt:&now,UserAgent:userAgent,CreatedAt:now,UpdatedAt:now}
 if err=o.sessions.Create(ctx,session);err!=nil{return auth.SessionResult{},err}
 return auth.SessionResult{Data:dto.LoginResult{AccountID:account.AccountID,PersonGID:account.PersonGID,SessionID:id,ExpiresIn:900,RequiresSetup:false},Cookies:auth.Cookies{Access:access,Refresh:refresh}},nil
}
