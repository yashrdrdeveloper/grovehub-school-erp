package impl

import (
 "context"
 "errors"
 "testing"
 "time"

 "grovehub/libs/schoolgateway"
 "grovehub/services/integration/internal/app/entity"
 mailcontracts "grovehub/services/integration/internal/app/providers/mailprovider/contracts"
 maildto "grovehub/services/integration/internal/app/providers/mailprovider/dto"
 connectionrepo "grovehub/services/integration/internal/app/repository/connection"
 helpers "grovehub/services/integration/internal/app/services/connection/helpers"
)

type tenantOAuthRepository struct {
 connectionrepo.ConnectionRepository
 state entity.OAuthState
 connection entity.Connection
 activated bool
}

func (r *tenantOAuthRepository) Get(_ context.Context, tenant, id string) (entity.Connection,error) {
 if tenant!=r.connection.TenantID || id!=r.connection.ID { return entity.Connection{},connectionrepo.ErrConnectionNotFound }
 return r.connection,nil
}
func (r *tenantOAuthRepository) GetByID(_ context.Context,id string)(entity.Connection,error){
 if id!=r.connection.ID{return entity.Connection{},connectionrepo.ErrConnectionNotFound}
 return r.connection,nil
}
func (r *tenantOAuthRepository) SaveOAuthState(_ context.Context,state entity.OAuthState)error{
 r.state=state;return nil
}
func (r *tenantOAuthRepository) ConsumeOAuthState(_ context.Context,hash string,now time.Time)(entity.OAuthState,error){
 if r.state.Hash!=hash||!r.state.ExpiresAt.After(now){return entity.OAuthState{},connectionrepo.ErrOAuthStateInvalid}
 state:=r.state;r.state=entity.OAuthState{};return state,nil
}
func (r *tenantOAuthRepository) Activate(context.Context,string,string,string,[]byte)error{
 r.activated=true;return nil
}

type tenantOAuthProvider struct {
 mailcontracts.OAuthStrategy
 exchanges int
 onExchange func()
}
func (p *tenantOAuthProvider) AuthorizationURL(state string)string{return state}
func (p *tenantOAuthProvider) Exchange(context.Context,string)(maildto.OAuthToken,error){
 p.exchanges++
 if p.onExchange!=nil{p.onExchange()}
 return maildto.OAuthToken{AccessToken:"test-access",RefreshToken:"test-refresh"},nil
}
func (p *tenantOAuthProvider) Mailbox(context.Context,maildto.OAuthToken)(string,error){return "mail@example.test",nil}

type tenantOAuthRegistry struct {
 mailcontracts.OAuthRegistry
 provider *tenantOAuthProvider
}
func (r tenantOAuthRegistry) Find(string)(mailcontracts.OAuthStrategy,error){return r.provider,nil}

func TestTenantOAuthGeneration(t *testing.T){
 for _,tc:=range []struct{
  name string
  legacy,disabled,reactivated,changeDuringExchange,unavailable bool
  exchanges int
  activate bool
 }{
  {name:"active",exchanges:1,activate:true},
  {name:"disabled",disabled:true},
  {name:"reactivated",reactivated:true},
  {name:"legacy",legacy:true},
  {name:"registry unavailable",unavailable:true},
  {name:"generation changes during provider exchange",changeDuringExchange:true,exchanges:1},
 }{
  t.Run(tc.name,func(t *testing.T){
   repo:=&tenantOAuthRepository{connection:entity.Connection{ID:"connection",TenantID:"school",Provider:"GMAIL",Status:entity.StatusPendingSetup}}
   provider:=&tenantOAuthProvider{}
   cipher,err:=helpers.NewTokenCipher(make([]byte,32));if err!=nil{t.Fatal(err)}
   operation:=New(repo,tenantOAuthRegistry{provider:provider},cipher)
   generation:=int64(1)
   active,available:=true,true
   ctx:=schoolgateway.WithTenantCheck(context.Background(),func(_ context.Context,tenant string)(int64,error){
    if tenant!="school"{t.Fatal("unexpected tenant",tenant)}
    if !active||!available{return 0,schoolgateway.ErrTenantDenied}
    return generation,nil
   })
   beginCtx:=ctx
   if tc.legacy{beginCtx=context.Background()}
   start,err:=operation.BeginOAuth(beginCtx,"school","connection");if err!=nil{t.Fatal(err)}
   want:=int64(1);if tc.legacy{want=0}
   if repo.state.SessionGeneration!=want{t.Fatal("wrong saved generation",repo.state.SessionGeneration)}
   if tc.disabled{active=false;generation++}
   if tc.reactivated{generation++}
   if tc.unavailable{available=false}
   if tc.changeDuringExchange{provider.onExchange=func(){generation++}}
   _,err=operation.FinishOAuth(ctx,"GMAIL","code",start.AuthorizationURL)
   if tc.activate&&err!=nil{t.Fatal(err)}
   if !tc.activate&&!errors.Is(err,helpers.ErrOAuthInvalid){t.Fatal("expected invalid state",err)}
   if provider.exchanges!=tc.exchanges||repo.activated!=tc.activate{t.Fatal("unexpected side effects",provider.exchanges,repo.activated)}
   if _,err=operation.FinishOAuth(ctx,"GMAIL","code",start.AuthorizationURL);!errors.Is(err,helpers.ErrOAuthInvalid){t.Fatal("replayed state accepted",err)}
  })
 }
}

func TestTenantOAuthBeginDeniesStaleActor(t *testing.T){
 repo:=&tenantOAuthRepository{connection:entity.Connection{ID:"connection",TenantID:"school",Provider:"GMAIL",Status:entity.StatusPendingSetup}}
 cipher,err:=helpers.NewTokenCipher(make([]byte,32));if err!=nil{t.Fatal(err)}
 operation:=New(repo,tenantOAuthRegistry{provider:&tenantOAuthProvider{}},cipher)
 ctx:=schoolgateway.WithTenantCheck(context.Background(),func(context.Context,string)(int64,error){return 2,nil})
 ctx=schoolgateway.WithPrincipal(ctx,schoolgateway.Principal{AccountID:"account",SessionID:"session",TenantID:"school",SessionGeneration:1})
 if _,err:=operation.BeginOAuth(ctx,"school","connection");!errors.Is(err,helpers.ErrOAuthInvalid){t.Fatal("stale actor accepted",err)}
 if repo.state.Hash!=""{t.Fatal("state issued for stale actor")}
}
