package refresh
import(
 "context"
 "errors"
 "time"
 dto "grovehub/services/identity/internal/app/dto/authentication"
 "grovehub/services/identity/internal/app/entity"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 sessionrepo "grovehub/services/identity/internal/app/repository/session"
 auth "grovehub/services/identity/internal/app/services/authentication"
 secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
)
type Operation struct{accounts accountrepo.Repository;sessions sessionrepo.Repository}
func New(a accountrepo.Repository,s sessionrepo.Repository)*Operation{return &Operation{accounts:a,sessions:s}}
func(o *Operation)Execute(ctx context.Context,raw string)(auth.RefreshResult,error){
 hash,err:=secrets.Digest(raw,48);if err!=nil{return auth.RefreshResult{},errors.New("invalid refresh token")}
 existing,err:=o.sessions.FindRefresh(ctx,hash,time.Now().UTC());if err!=nil{return auth.RefreshResult{},errors.New("invalid refresh token")}
 account,err:=o.accounts.FindByID(ctx,existing.AccountID);if err!=nil||account.Status!=entity.Active{return auth.RefreshResult{},errors.New("invalid refresh token")}
 access,err:=secrets.NewSecret(32);if err!=nil{return auth.RefreshResult{},err}
 refresh,err:=secrets.NewSecret(48);if err!=nil{return auth.RefreshResult{},err}
 now:=time.Now().UTC();accessHash,_:=secrets.Digest(access,32);refreshHash,_:=secrets.Digest(refresh,48)
 updated:=existing;updated.AccessHash=accessHash;updated.RefreshHash=refreshHash;updated.AccessExpiresAt=now.Add(15*time.Minute);updated.UpdatedAt=now
 if err=o.sessions.Rotate(ctx,updated,hash);err!=nil{return auth.RefreshResult{},errors.New("invalid refresh token")}
 return auth.RefreshResult{Data:dto.RefreshResult{SessionID:existing.SessionID,ExpiresIn:900,Refreshed:true},Cookies:auth.Cookies{Access:access,Refresh:refresh}},nil
}
